package client

import (
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

type Client struct {
	client  *uhttp.BaseHttpClient
	baseURL string
}

func NewClient(client *uhttp.BaseHttpClient, baseURL string) *Client {
	return &Client{
		client:  client,
		baseURL: baseURL,
	}
}
