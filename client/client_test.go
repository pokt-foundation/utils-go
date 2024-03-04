package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// Initialize a default test client with a 5 second timeout
var testClient = NewDefaultClient()

/* ---------- Unit Tests (make test server API calls) ---------- */

type Response struct {
	Message string `json:"message"`
}

func Test_Get_Unit(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		want           Response
		expectError    bool
	}{
		{
			name:           "should successfully retrieve data",
			serverResponse: `{"message":"success"}`,
			want:           Response{Message: "success"},
			expectError:    false,
		},
		{
			name:           "should handle server error",
			serverResponse: `Server Error`,
			want:           Response{},
			expectError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.expectError {
					http.Error(w, test.serverResponse, http.StatusInternalServerError)
				} else {
					fmt.Fprint(w, test.serverResponse)
				}
			}))
			defer server.Close()

			// Execute Get request
			got, err := Get[Response](GetReqConfig{
				URL:    server.URL,
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(test.want, got)
			}
		})
	}
}

func Test_Post_Unit(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		serverResponse string
		want           Response
		expectError    bool
	}{
		{
			name:           "should successfully post data",
			body:           map[string]string{"input": "test"},
			serverResponse: `{"message":"success"}`,
			want:           Response{Message: "success"},
			expectError:    false,
		},
		{
			name:           "should handle server error on post",
			body:           map[string]string{"input": "test"},
			serverResponse: `Server Error`,
			want:           Response{},
			expectError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.expectError {
					http.Error(w, test.serverResponse, http.StatusInternalServerError)
				} else {
					fmt.Fprint(w, test.serverResponse)
				}
			}))
			defer server.Close()

			// Execute Post request
			got, err := Post[Response](PostReqConfig{
				URL:    server.URL,
				Body:   test.body,
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(test.want, got)
			}
		})
	}
}

func Test_Put_Unit(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		serverResponse string
		want           Response
		expectError    bool
	}{
		{
			name:           "should successfully put data",
			body:           map[string]string{"input": "update"},
			serverResponse: `{"message":"updated"}`,
			want:           Response{Message: "updated"},
			expectError:    false,
		},
		{
			name:           "should handle server error on put",
			body:           map[string]string{"input": "update"},
			serverResponse: `Server Error`,
			want:           Response{},
			expectError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.expectError {
					http.Error(w, test.serverResponse, http.StatusInternalServerError)
				} else {
					fmt.Fprint(w, test.serverResponse)
				}
			}))
			defer server.Close()

			// Execute Put request
			got, err := Put[Response](PutReqConfig{
				URL:    server.URL,
				Body:   test.body,
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(test.want, got)
			}
		})
	}
}

func Test_Patch_Unit(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		serverResponse string
		want           Response
		expectError    bool
	}{
		{
			name:           "should successfully patch data",
			body:           map[string]string{"input": "modify"},
			serverResponse: `{"message":"modified"}`,
			want:           Response{Message: "modified"},
			expectError:    false,
		},
		{
			name:           "should handle server error on patch",
			body:           map[string]string{"input": "modify"},
			serverResponse: `Server Error`,
			want:           Response{},
			expectError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.expectError {
					http.Error(w, test.serverResponse, http.StatusInternalServerError)
				} else {
					fmt.Fprint(w, test.serverResponse)
				}
			}))
			defer server.Close()

			// Execute Patch request
			got, err := Patch[Response](PatchReqConfig{
				URL:    server.URL,
				Body:   test.body,
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(test.want, got)
			}
		})
	}
}

func Test_Delete_Unit(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		wantStatusCode int
		expectError    bool
	}{
		{
			name:           "should successfully delete data",
			serverResponse: `{"message":"deleted"}`,
			wantStatusCode: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "should handle server error on delete",
			serverResponse: `Server Error`,
			wantStatusCode: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.expectError {
					http.Error(w, test.serverResponse, test.wantStatusCode)
				} else {
					w.WriteHeader(test.wantStatusCode)
					fmt.Fprint(w, test.serverResponse)
				}
			}))
			defer server.Close()

			// Execute Delete request
			config := DeleteReqConfig{
				URL:    server.URL,
				Client: testClient,
			}
			_, err := Delete[Response](config)

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(http.StatusOK, test.wantStatusCode)
			}
		})
	}
}

func Test_parseErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		wantErrMessage string
	}{
		{
			name: "should return basic error for non-JSON response",
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     "500 Internal Server Error",
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			wantErrMessage: "response not OK. 500 Internal Server Error",
		},
		{
			name: "should return detailed error for JSON response with common error key",
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     "400 Bad Request",
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Invalid request"}`)),
			},
			wantErrMessage: "response not OK. 400 Bad Request: Invalid request",
		},
		{
			name: "should handle empty error message in JSON response",
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     "400 Bad Request",
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": ""}`)),
			},
			wantErrMessage: "response not OK. 400 Bad Request",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := parseErrorResponse(test.response)
			require.EqualError(t, err, test.wantErrMessage)
		})
	}
}

/* ---------- Integration Tests (make actual API calls) ---------- */

type (
	PokemonAPIResponse struct {
		Name string `json:"name"`
	}

	HTTPBinResponse struct {
		JSON string `json:"json"`
	}
)

func Test_Get_Integration_Unit(t *testing.T) {
	tests := []struct {
		name        string
		pokemonID   string
		wantName    string
		expectError bool
	}{
		{
			name:        "should successfully retrieve bulbasaur name",
			pokemonID:   "1",
			wantName:    "bulbasaur",
			expectError: false,
		},
		{
			name:        "should successfully retrieve scyther name",
			pokemonID:   "123",
			wantName:    "scyther",
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Execute Get request
			resp, err := Get[PokemonAPIResponse](GetReqConfig{
				URL:    fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", test.pokemonID),
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(test.wantName, resp.Name)
			}
		})
	}
}

func Test_Post_Integration(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		wantJSON       HTTPBinResponse
		wantStatusCode int
		expectError    bool
	}{
		{
			name:           "should successfully post data to httpbin and validate json response",
			body:           `{"lisa": "needs braces"}`,
			wantJSON:       HTTPBinResponse{JSON: `{"lisa": "needs braces"}`},
			wantStatusCode: 200,
			expectError:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Execute Post request
			resp, err := Post[HTTPBinResponse](PostReqConfig{
				URL:    "https://httpbin.org/post",
				Body:   test.body,
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(test.wantJSON, resp)
			}
		})
	}
}

func Test_Put_Integration(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		wantJSON       HTTPBinResponse
		wantStatusCode int
		expectError    bool
	}{
		{
			name:           "should successfully put data to httpbin and validate json response",
			body:           `{"quote":"I am that guy"}`,
			wantJSON:       HTTPBinResponse{JSON: `{"quote":"I am that guy"}`},
			wantStatusCode: 200,
			expectError:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Execute Put request
			resp, err := Put[HTTPBinResponse](PutReqConfig{
				URL:    "https://httpbin.org/put",
				Body:   test.body,
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(test.wantJSON, resp)
			}
		})
	}
}

func Test_Patch_Integration(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		wantJSON       HTTPBinResponse
		wantStatusCode int
		expectError    bool
	}{
		{
			name:           "should successfully patch data to httpbin and validate json response",
			body:           `{"quote":"Even the smallest person can change the course of the future."}`,
			wantJSON:       HTTPBinResponse{JSON: `{"quote":"Even the smallest person can change the course of the future."}`},
			wantStatusCode: 200,
			expectError:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Execute Patch request
			resp, err := Patch[HTTPBinResponse](PatchReqConfig{
				URL:    "https://httpbin.org/patch",
				Body:   test.body,
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.Equal(test.wantJSON, resp)
			}
		})
	}
}

func Test_Delete_Integration(t *testing.T) {
	tests := []struct {
		name           string
		wantStatusCode int
		expectError    bool
	}{
		{
			name:           "should successfully delete data at httpbin and validate json response",
			wantStatusCode: 200,
			expectError:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := require.New(t)

			// Execute Delete request
			_, err := Delete[HTTPBinResponse](DeleteReqConfig{
				URL:    "https://httpbin.org/delete",
				Client: testClient,
			})

			// Validate
			if test.expectError {
				c.Error(err)
			} else {
				c.NoError(err)
			}
		})
	}
}
