package client

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

type Client struct {
	client    *uhttp.BaseHttpClient
	baseURL   string
	authToken string
}

func NewClient(client *uhttp.BaseHttpClient, baseURL string, authToken string) *Client {
	return &Client{
		client:    client,
		baseURL:   baseURL,
		authToken: authToken,
	}
}

func (c *Client) do(ctx context.Context, method, pathAndQuery string, body io.Reader) (*http.Response, error) {
	fullUrl := &url.URL{
		Path:   pathAndQuery,
		Host:   c.baseURL,
		Scheme: "http",
	}
	parsedUrl, err := url.QueryUnescape(fullUrl.String())
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, parsedUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	return c.client.Do(req)
}
