package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type queryHTTP struct {
	name    string
	url     string
	method  string
	headers map[string]string
	body    string
	// nolint:unused
	dependencyName string
}

func scan(r io.Reader) ([]queryHTTP, error) {
	if r == nil {
		return nil, nil
	}

	s := bufio.NewScanner(r)

	result := make([]queryHTTP, 0)

	var (
		current  *queryHTTP
		openBody bool
	)

	for s.Scan() {
		line := s.Text()

		// parse start
		if isName(line) || (current == nil && isURL(line)) {
			openBody = false
			if current != nil {
				result = append(result, *current)
			}

			current = &queryHTTP{
				name:    parseName(line, len(result)+1),
				headers: make(map[string]string),
			}

			if isURL(line) {
				current.url = parseURL(line)
				current.method = parseMethod(line)
			}

			continue
		}

		if current == nil { // skip
			continue
		}

		if line == "" && current.method == "POST" { // empty line
			openBody = true
			continue
		}

		if openBody {
			current.body += strings.TrimSpace(line)
			continue
		}

		if isURL(line) {
			current.url = parseURL(line)
			current.method = parseMethod(line)
			continue
		}

		if isHeader(line, current) {
			continue
		}

	}

	if current != nil {
		result = append(result, *current)
	}

	return result, nil
}

func isURL(line string) bool {
	if len(line) < 4 {
		return false
	}

	if line[:4] == "http" {
		return true
	}
	if line[:3] == "GET" {
		return true
	}

	if line[:4] == "POST" {
		return true
	}

	return false
}

func parseName(line string, numberInFile int) string {
	if isURL(line) {
		return fmt.Sprintf("###query #%d", numberInFile)
	}

	return line
}

func parseMethod(line string) string {
	if line[:4] == "POST" {
		return "POST"
	}

	return "GET"
}

func parseURL(line string) string {
	if line[:3] == "GET" {
		return strings.TrimSpace(line[3:])
	}

	if line[:4] == "POST" {
		return strings.TrimSpace(line[4:])
	}

	return line
}

func isName(line string) bool {
	if len(line) < 4 {
		return false
	}

	if line[:3] == "###" {
		return true
	}

	return false
}

func isHeader(line string, query *queryHTTP) bool {
	if s := strings.Split(line, ":"); len(s) == 2 {
		k, v := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
		if len(k) > 0 && len(v) > 0 {
			query.headers[k] = v

			return true
		}

	}

	return false
}
