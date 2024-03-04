// package client provides utility methods for making all HTTP requests
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
	defaultRetries = 0
)

var ErrResponseNotOK error = errors.New("response not OK")

type (
	// GetReqConfig is a struct that holds the configuration for a GET request.
	// If a Client is not provided, a default client with a timeout of 5 seconds will be used.
	GetReqConfig struct {
		URL     string
		Headers map[string]string
		Client  *http.Client
	}

	// PostReqConfig is a struct that holds the configuration for a POST request.
	// If a Client is not provided, a default client with a timeout of 5 seconds will be used.
	PostReqConfig struct {
		URL     string
		Body    any
		Headers map[string]string
		Client  *http.Client
	}

	// PutReqConfig is a struct that holds the configuration for a PUT request.
	// If a Client is not provided, a default client with a timeout of 5 seconds will be used.
	PutReqConfig struct {
		URL     string
		Body    any
		Headers map[string]string
		Client  *http.Client
	}

	// PatchReqConfig is a struct that holds the configuration for a PATCH request.
	// If a Client is not provided, a default client with a timeout of 5 seconds will be used.
	PatchReqConfig struct {
		URL     string
		Body    any
		Headers map[string]string
		Client  *http.Client
	}
	// DeleteReqConfig is a struct that holds the configuration for a DELETE request.
	// If a Client is not provided, a default client with a timeout of 5 seconds will be used.
	DeleteReqConfig struct {
		URL     string
		Headers map[string]string
		Client  *http.Client
	}

	// retryTransport is a custom transport for handling retries
	retryTransport struct {
		underlying http.RoundTripper
		retries    int
	}
)

// NewDefaultClient returns Client struct with default timeout of 5 seconds
func NewDefaultClient() *http.Client {
	return &http.Client{Timeout: defaultTimeout}
}

// NewCustomClient returns Client struct with custom timeout and retries
func NewCustomClient(timeout time.Duration, retries int) *http.Client {
	client := &http.Client{Timeout: timeout}

	if retries > 0 {
		client.Transport = &retryTransport{
			underlying: client.Transport,
			retries:    retries,
		}
	}

	return client
}

// Get makes a GET request to the given URL with the given headers and client
func Get[R any](c GetReqConfig) (R, error) {
	if c.Client == nil {
		c.Client = NewDefaultClient()
	}

	var data R

	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if err != nil {
		return data, err
	}

	for header, headerValue := range c.Headers {
		req.Header.Set(header, headerValue)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, err
}

// Post makes a POST request to the given URL with the given body, headers and client
func Post[R any](c PostReqConfig) (R, error) {
	if c.Client == nil {
		c.Client = NewDefaultClient()
	}

	var data R

	postData, err := json.Marshal(c.Body)
	if err != nil {
		return data, err
	}

	req, err := http.NewRequest(http.MethodPost, c.URL, bytes.NewBuffer(postData))
	if err != nil {
		return data, err
	}

	req.Header.Set("Content-Type", "application/json")
	for header, headerValue := range c.Headers {
		req.Header.Set(header, headerValue)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data, parseErrorResponse(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, err
}

// Put makes a PUT request to the given URL with the given body, headers and client
func Put[R any](c PutReqConfig) (R, error) {
	if c.Client == nil {
		c.Client = NewDefaultClient()
	}

	var data R

	putData, err := json.Marshal(c.Body)
	if err != nil {
		return data, err
	}

	req, err := http.NewRequest(http.MethodPut, c.URL, bytes.NewBuffer(putData))
	if err != nil {
		return data, err
	}

	req.Header.Set("Content-Type", "application/json")
	for header, headerValue := range c.Headers {
		req.Header.Set(header, headerValue)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data, parseErrorResponse(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, err
}

// Patch makes a PATCH request to the given URL with the given body, headers and client
func Patch[R any](c PatchReqConfig) (R, error) {
	if c.Client == nil {
		c.Client = NewDefaultClient()
	}

	var data R

	patchData, err := json.Marshal(c.Body)
	if err != nil {
		return data, err
	}

	req, err := http.NewRequest(http.MethodPatch, c.URL, bytes.NewBuffer(patchData))
	if err != nil {
		return data, err
	}

	req.Header.Set("Content-Type", "application/json")
	for header, headerValue := range c.Headers {
		req.Header.Set(header, headerValue)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data, parseErrorResponse(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, err
}

// Delete makes a DELETE request to the given URL with the given headers and client
func Delete[R any](c DeleteReqConfig) (R, error) {
	if c.Client == nil {
		c.Client = NewDefaultClient()
	}

	var data R

	req, err := http.NewRequest(http.MethodDelete, c.URL, nil)
	if err != nil {
		return data, err
	}

	for header, headerValue := range c.Headers {
		req.Header.Set(header, headerValue)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data, parseErrorResponse(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, err
}

// retryTransport is a custom transport for handling retries
func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.underlying
	if rt == nil {
		rt = http.DefaultTransport
	}

	var resp *http.Response
	var err error

	// Cache request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i <= t.retries; i++ {
		// Recreate body reader
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		resp, err = rt.RoundTrip(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}

		if i < t.retries {
			time.Sleep(time.Duration(i*i) * 100 * time.Millisecond)
		}
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// commonErrKeys is a list of common error message keys
var commonErrKeys = []string{
	"error", "message", "errorMessage",
	"error_description", "errorDescription",
	"error_message", "errorMsg",
	"err", "errorText", "error_text",
	"description", "desc", "errorDetails", "error_details",
}

// parseErrorResponse is a helper function to parse error responses from the server
func parseErrorResponse(errResponse *http.Response) error {
	code := errResponse.StatusCode
	text := http.StatusText(code)
	basicError := fmt.Errorf("%s. %d %s", ErrResponseNotOK, code, text)

	body, err := io.ReadAll(errResponse.Body)
	if err != nil {
		return basicError
	}

	var errorObj map[string]interface{}
	if err := json.Unmarshal(body, &errorObj); err == nil {
		for _, key := range commonErrKeys {
			if msg, ok := errorObj[key].(string); ok && msg != "" {
				return fmt.Errorf("%s: %s", basicError, msg)
			}
		}
	}

	return basicError
}
