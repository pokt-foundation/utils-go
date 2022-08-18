package client

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/utils-go/mock-client"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultClient(t *testing.T) {
	c := require.New(t)

	client := NewDefaultClient()
	c.NotEmpty(client)

	client = NewCustomClient(5, 3*time.Second)
	c.NotEmpty(client)
}

func TestClient_PostWithURLJSONParams(t *testing.T) {
	c := require.New(t)

	client := NewDefaultClient()
	c.NotEmpty(client)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock.AddMockedResponseFromFile(http.MethodPost, "https://dummy.com", http.StatusCreated, "../mock-client/samples/dummy.json")

	response, err := client.PostWithURLJSONParams("https://dummy.com", map[string]string{
		"ohana": "family",
	}, http.Header{})
	c.NoError(err)

	c.NotEmpty(response)
	c.Equal(http.StatusCreated, response.StatusCode)
	c.NoError(response.Body.Close())
}

func TestClient_GetWithURLAndParams(t *testing.T) {
	c := require.New(t)

	client := NewDefaultClient()
	c.NotEmpty(client)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock.AddMockedResponseFromFile(http.MethodGet, "https://dummy.com", http.StatusOK, "../mock-client/samples/dummy.json")

	params := url.Values{}
	params.Set("family", "ohana")

	response, err := client.GetWithURLAndParams("https://dummy.com", params, http.Header{})
	c.NoError(err)

	c.NotEmpty(response)
	c.Equal(http.StatusOK, response.StatusCode)
	c.NoError(response.Body.Close())
}
