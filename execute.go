package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func executeQuery(queryName string, queryNumber int, result ...queryHTTP) {
	for i := range result {
		q := result[i]

		if queryName != "" && queryName != q.name {
			continue
		}
		if queryNumber != 0 && queryNumber != q.queryNumber {
			continue
		}

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
