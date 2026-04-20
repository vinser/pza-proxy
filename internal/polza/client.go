package polza

import (
	"bytes"
	"net/http"
	"time"
)

const polzaURL = "https://polza.ai/api/v1/chat/completions"

// Client is a thin Polza HTTP client.
type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 0, // streaming: no global timeout
		},
	}
}

func (c *Client) Do(body []byte, stream bool) (*http.Response, error) {
	req, err := http.NewRequest("POST", polzaURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey) // already contains "Bearer ..."

	if !stream {
		c.http.Timeout = 60 * time.Second
	}

	return c.http.Do(req)
}
