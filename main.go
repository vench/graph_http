package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	version string
	build   string
)

func main() {
	envPath := flag.String("env", ".env", "Path to environment config file")
	inputFile := flag.String("path", "example.http", "Path to file for parsing .http format")
	showVersion := flag.Bool("version", false, "Show current version")

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

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("failed to open file[%s] : %v", *inputFile, err)
	}
	defer f.Close()

	result, err := scan(f)
	if err != nil {
		log.Fatalf("failed to scan file[%s] : %v", *inputFile, err)
	}

	for i := range result {
		q := result[i]

		fmt.Printf("run query : %s\n\n", q.name)

		request, err := http.NewRequest(q.method, q.url, bytes.NewBufferString(q.body))
		if err != nil {
			log.Fatalf("failed to create new http request: %v", err)
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalf("failed to do http request: %v", err)
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("failed to read body: %v", err)
		}

		for k := range response.Header {
			fmt.Printf("%s: %s\n", k, strings.Join(response.Header[k], ","))
		}

		fmt.Printf("\n%s\n\n", string(data))
	}
}
