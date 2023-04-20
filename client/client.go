// Package client is a wrapper for the default go client with default and configurable options
// It also has some utility functions to do actual requests
package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultHTTPClientTimeout = 5 * time.Second
	defaultHTTPClientRetries = 0
	backoffTime              = 2 * time.Second
)

// Client is a wrapper for the default client
type Client struct {
	Client  *http.Client
	retries int
}

// NewDefaultClient returns httpclient instance with default config
func NewDefaultClient() *Client {
	return &Client{
		Client: &http.Client{
			Timeout: defaultHTTPClientTimeout,
		},
		retries: defaultHTTPClientRetries,
	}
}

// NewCustomClient returns httpclient instance with given custom config
func NewCustomClient(retries int, timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
		},
		retries: retries,
	}
}

func getJSONBodyFromParams(params any) (io.Reader, error) {
	if params == nil {
		return nil, nil
	}

	rawBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return bytes.NewBufferString(string(rawBody)), nil
}

// PostWithURLJSONParams does post request with JSON param
func (c *Client) PostWithURLJSONParams(url string, params any, headers http.Header) (*http.Response, error) {
	body, err := getJSONBodyFromParams(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	headers.Set("Content-Type", "application/json")
	headers.Set("Connection", "close")

	req.Header = headers

	return c.DoRequestWithRetries(req)
}

// PostWithURLEncodedParams does post request with URL encoded params
func (c *Client) PostWithURLEncodedParams(url string, params url.Values, headers http.Header) (*http.Response, error) {
	var body io.Reader

	if len(params) != 0 {
		body = strings.NewReader(params.Encode())
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Connection", "close")

	req.Header = headers

	return c.DoRequestWithRetries(req)
}

// PutWithURLJSONParams does post request with JSON param
func (c *Client) PutWithURLJSONParams(url string, params any, headers http.Header) (*http.Response, error) {
	body, err := getJSONBodyFromParams(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	headers.Set("Content-Type", "application/json")
	headers.Set("Connection", "close")

	req.Header = headers

	return c.DoRequestWithRetries(req)
}

// GetWithURLAndParams does get request with url values as params
func (c *Client) GetWithURLAndParams(rawURL string, params url.Values, headers http.Header) (*http.Response, error) {
	urlStruct, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	urlStruct.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, urlStruct.String(), nil)
	if err != nil {
		return nil, err
	}

	return c.DoRequestWithRetries(req)
}

// DoRequestWithRetries does requests with the retries set on client and backoff time
// just retries request with status code 5xx
func (c *Client) DoRequestWithRetries(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// at least one attempt is made, regardless of how many retries were on config
	attempts := c.retries + 1

	for i := 0; i < attempts; i++ {
		resp, err = c.Client.Do(req)
		if err != nil {
			return nil, err
		}

		// Just retry 5xx responses
		if string(resp.Status[0]) != "5" {
			break
		}

		// On the last attempt there's no reason to wait the backoff time
		if i != attempts-1 {
			time.Sleep(backoffTime)
		}
	}

	return resp, nil
}
