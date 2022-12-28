package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_scan(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		input  io.Reader
		result []queryHTTP
		err    error
	}{
		{
			name: "empty",
		},
		{
			name:  "ok - one line",
			input: bytes.NewBufferString("http://localhost"),
			result: []queryHTTP{
				{
					name:        parseName("http://localhost", 1),
					queryNumber: 1,
					url:         "http://localhost",
					method:      http.MethodGet,
					headers:     make(map[string]string),
				},
			},
		},
		{
			name: "ok - multi line",
			input: bytes.NewBufferString(`
http://localhost
Content-Type: application/json
X-Requested-With: XMLHttpRequest

###some name
GET http://localhost:8080/?foo=bar

## skip line

### post query < some name
POST https://test.local/foo

{
  "json":1,
  "date":"2022-11-12"
}

`),
			result: []queryHTTP{
				{
					name:        "###query #1",
					queryNumber: 1,
					url:         "http://localhost",
					method:      http.MethodGet,
					headers: map[string]string{
						"Content-Type":     "application/json",
						"X-Requested-With": "XMLHttpRequest",
					},
				},
				{
					name:        "###some name",
					queryNumber: 2,
					url:         "http://localhost:8080/?foo=bar",
					method:      "GET",
					headers:     make(map[string]string),
				},
				{
					name:           "### post query",
					dependencyName: "### some name",
					queryNumber:    3,
					url:            "https://test.local/foo",
					method:         "POST",
					headers:        make(map[string]string),
					body:           "{\"json\":1,\"date\":\"2022-11-12\"}",
				},
			},
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := scan(tc.input)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}
