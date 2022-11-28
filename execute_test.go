package main

import "testing"

func Test_execute(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		queryName   string
		queryNumber int
		queries     []queryHTTP
	}{
		{
			name: "empty",
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			executeQuery(tc.queryName, tc.queryNumber, tc.queries...)
		})
	}
}
