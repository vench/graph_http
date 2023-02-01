package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	postMethod = "POST"
	getMethod  = "GET"
)

type queryHTTP struct {
	name        string
	queryNumber int

	url     string
	method  string
	headers map[string]string
	body    string

	dependencyName string
	output         string
}

var errEmptyReader = errors.New("empty reader")

func scan(r io.Reader) ([]queryHTTP, error) {
	if r == nil {
		return nil, errEmptyReader
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

			pName := parseName(line, len(result)+1)
			name, dependencyName := splitDependencyName(pName)
			current = &queryHTTP{
				name:           name,
				dependencyName: dependencyName,
				headers:        make(map[string]string),
				queryNumber:    len(result) + 1,
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

		if line == "" && current.method == postMethod { // empty line
			openBody = !openBody

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

		if isOutput(line, current) {
			continue
		}

		if isProcessed(line) {
			_ = line

			continue
		}
	}

	if current != nil {
		result = append(result, *current)
	}

	return result, nil
}

func isProcessed(line string) bool {
	if line != "" && line[:1] == ">" {
		return true
	}

	return false
}

func isOutput(line string, query *queryHTTP) bool {
	if line != "" && line[:2] == ">>" {
		query.output = strings.TrimSpace(line[2:])

		return true
	}

	return false
}

func isURL(line string) bool {
	if len(line) < 4 {
		return false
	}

	if line[:4] == "http" {
		return true
	}
	if line[:3] == getMethod {
		return true
	}

	if line[:4] == postMethod {
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

func splitDependencyName(name string) (string, string) {
	if s := strings.Split(name, "<"); len(s) == 2 {
		if !isName(s[1]) && len(s[1]) != 0 {
			return strings.TrimSpace(s[0]), "### " + strings.TrimSpace(s[1])
		}

		return strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
	}

	return name, ""
}

func parseMethod(line string) string {
	if line[:4] == postMethod {
		return postMethod
	}

	return getMethod
}

func parseURL(line string) string {
	if line[:3] == getMethod {
		return strings.TrimSpace(line[3:])
	}

	if line[:4] == postMethod {
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
