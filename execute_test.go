package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_execute(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		queries []queryHTTP
	}{
		{
			name: "empty",
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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
