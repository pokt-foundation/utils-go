// Package client is a wrapper for heimdall github.com/gojektech/heimdall with default and configurable options
// It also has some utility functions to do actual requests
package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
)

const (
	defaultHTTPClientTimeout = 5 * time.Second
	defaultHTTPClientRetries = 0

	initialBackoffTimeout = 2 * time.Millisecond
	maxBackoffTimeout     = 9 * time.Millisecond
	exponentFactor        = 2
	maxJitterInterval     = 2 * time.Millisecond
)

var (
	backoff = heimdall.NewExponentialBackoff(initialBackoffTimeout, maxBackoffTimeout, exponentFactor, maxJitterInterval)
	retrier = heimdall.NewRetrier(backoff)
)

// Client is a wrapper for the heimdall client
type Client struct {
	*httpclient.Client
}

// NewDefaultClient returns httpclient instance with default config
func NewDefaultClient() *Client {
	return &Client{
		Client: httpclient.NewClient(
			httpclient.WithHTTPTimeout(defaultHTTPClientTimeout),
			httpclient.WithRetryCount(defaultHTTPClientRetries),
			httpclient.WithRetrier(retrier),
		),
	}
}

// NewCustomClient returns httpclient instance with given custom config
func NewCustomClient(retries int, timeout time.Duration) *Client {
	return &Client{
		Client: httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetryCount(retries),
			httpclient.WithRetrier(retrier),
		),
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

	headers.Set("Content-Type", "application/json")
	headers.Set("Connection", "close")

	return c.Post(url, body, headers)
}

// PutWithURLJSONParams does post request with JSON param
func (c *Client) PutWithURLJSONParams(url string, params any, headers http.Header) (*http.Response, error) {
	body, err := getJSONBodyFromParams(params)
	if err != nil {
		return nil, err
	}

	headers.Set("Content-Type", "application/json")
	headers.Set("Connection", "close")

	return c.Put(url, body, headers)
}

// GetWithURLAndParams does get request with url values as params
func (c *Client) GetWithURLAndParams(rawURL string, params url.Values, headers http.Header) (*http.Response, error) {
	urlStruct, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	urlStruct.RawQuery = params.Encode()

	return c.Get(urlStruct.String(), headers)
}
