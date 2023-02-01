package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_execute(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		handler func(writer http.ResponseWriter, request *http.Request)
		queries []queryHTTP
	}{
		{
			name: "empty",
		},
		{
			name: "ok",
			queries: []queryHTTP{
				{
					name: "query1",
				},
				{
					name:           "query2",
					dependencyName: "query1",
				},
			},
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.handler == nil {
				tc.handler = func(writer http.ResponseWriter, request *http.Request) {
					t.Log("x")
					writer.WriteHeader(http.StatusOK)
				}
			}

			serv := httptest.NewServer(http.HandlerFunc(tc.handler))
			defer serv.Close()

			for i := range tc.queries {
				tc.queries[i].url = serv.URL
			}

			executeQuery(tc.queries...)
		})
	}
}

func Test_filterQueries(t *testing.T) {
	t.Parallel()

	queries := []queryHTTP{
		{
			name:        "test1",
			queryNumber: 1,
		},
		{
			name:        "test2",
			queryNumber: 2,
		},
		{
			name:        "test3",
			queryNumber: 3,
		},
	}

	tt := []struct {
		name        string
		queryName   string
		queryNumber int
		queries     []queryHTTP
		out         []queryHTTP
	}{
		{
			name: "empty",
			out:  queries[:0],
		},
		{
			name:    "with out filters",
			queries: queries,
			out:     queries,
		},
		{
			name:      "by name",
			queryName: "test2",
			queries:   queries,
			out:       queries[1:2],
		},
		{
			name:        "by number",
			queryNumber: 3,
			queries:     queries,
			out:         queries[2:],
		},
		{
			name:        "name priority",
			queryName:   "test2",
			queryNumber: 1,
			queries:     queries,
			out:         queries[1:2],
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := filterQueries(tc.queryName, tc.queryNumber, tc.queries)
			require.Equal(t, tc.out, out)
		})
	}
}
