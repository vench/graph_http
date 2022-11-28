package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	version string
	build   string
)

func main() {
	var (
		envPath     = flag.String("env", ".env", "Path to environment config file")
		inputFile   = flag.String("path", "example.http", "Path to file for parsing .http format")
		showVersion = flag.Bool("version", false, "Show current version")

		queryName   string
		queryURL    string
		queryNumber int
	)

	flag.StringVar(&queryURL, "query_url", "", "Direct URL for run")
	flag.StringVar(&queryURL, "u", "", "Alias for -query_url")
	flag.StringVar(&queryName, "query_name", "", "Filter by query name")
	flag.StringVar(&queryName, "q", "", "Alias for -query_url")
	flag.IntVar(&queryNumber, "query_number", 0, "Filter by query number position in file")
	flag.IntVar(&queryNumber, "n", 0, "Alias for -query_number")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s, build: %s\n", version, build)
		return
	}

	if _, err := os.Stat(*envPath); err == nil {
		if err = godotenv.Load(*envPath); err != nil {
			log.Fatalf("failed to loading .env file: %v", err)
		}
	}

	if queryURL != "" {
		executeQuery("", 0, queryHTTP{
			name:   queryURL,
			url:    parseURL(queryURL),
			method: parseMethod(queryURL),
		})

		return
	}

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("failed to open file[%s] : %v", *inputFile, err)
	}
	defer f.Close()

	result, err := scan(f)
	if err != nil {
		log.Fatalf("failed to scan file[%s] : %v", *inputFile, err)
	}

	executeQuery(queryName, queryNumber, result...)
}
