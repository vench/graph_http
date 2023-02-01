package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultTimeOut = 5 * time.Second

func executeQuery(result ...queryHTTP) {
	var wg sync.WaitGroup

	dm := make(map[string]chan struct{})
	for i := range result {
		dm[result[i].name] = make(chan struct{})
	}

	for i := range result {
		q := result[i]

		wg.Add(1)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			defer close(dm[q.name])

			if c, ok := dm[q.dependencyName]; ok {
				// wait dependency query
				log.Printf("query %s wait %s\n", q.name, q.dependencyName)
				<-c
			}

			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
			defer cancel()

			log.Printf("begin run query : %s\n", q.name)

			request, err := http.NewRequestWithContext(ctx, q.method, q.url, bytes.NewBufferString(q.body))
			if err != nil {
				log.Fatalf("failed to create new http request: %v", err)
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				log.Fatalf("failed to do http request[%s]: %v", q.name, err)
			}
			defer response.Body.Close()

			log.Printf("done run query : %s, status code: %d\n", q.name, response.StatusCode)

			data, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatalf("failed to read body: %v", err)
			}

			outWriter := os.Stdout
			if q.output != "" {
				f, err := os.Open(q.output)
				if err != nil {
					log.Fatalf("failed to open file %s: %v", q.output, err)
				}
				defer outWriter.Close()

				outWriter = f
			}

			for k := range response.Header {
				fmt.Fprintf(outWriter, "%s: %s\n", k, strings.Join(response.Header[k], ","))
			}

			fmt.Fprintf(outWriter, "\n%s\n\n", string(data))
		}(&wg)
	}

	wg.Wait()
}

func filterQueries(queryName string, queryNumber int, queries []queryHTTP) []queryHTTP {
	result := make([]queryHTTP, 0, len(queries))

	for i := range queries {
		q := queries[i]

		if queryName != "" && queryName != q.name {
			continue
		}

		if queryName == "" && queryNumber != 0 && queryNumber != q.queryNumber {
			continue
		}

		result = append(result, q)
	}

	return result
}
