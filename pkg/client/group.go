package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type ContextItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Node struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value bool   `json:"value"`
	//Context []map[string]string `json:"context"`
}

type Group struct {
	Name        string         `json:"name"`
	DisplayName string         `json:"displayName"`
	Nodes       []*Node        `json:"nodes"`
	Metadata    map[string]any `json:"metadata"`
}

func (o *Client) GetGroup(ctx context.Context, id string) (*Group, error) {
	req, err := o.client.NewRequest(ctx, "GET", &url.URL{
		Path:   fmt.Sprintf("group/%s", id),
		Host:   o.baseURL,
		Scheme: "http",
	})
	if err != nil {
		return nil, err
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if err := resp.Body.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}
	g := Group{}
	if err := json.Unmarshal(body, &g); err != nil {
		if err := resp.Body.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	if g.DisplayName == "" {
		g.DisplayName = g.Name
	}

	return &g, nil
}

func (o *Client) ListAllGroups(ctx context.Context) ([]*Group, error) {
	req, err := o.client.NewRequest(ctx, "GET", &url.URL{
		Path:   "group",
		Host:   o.baseURL,
		Scheme: "http",
	})
	if err != nil {
		return nil, err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	groups := []string{}
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}

	resources := []*Group{}
	for _, group := range groups {
		g, err := o.GetGroup(ctx, group)
		if err != nil {
			return nil, err
		}

		resources = append(resources, g)
	}
	return resources, nil
}
