package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type ContextItem struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Node struct {
	Key     string `json:"key,omitempty"`
	Type    string `json:"type,omitempty"`
	Value   bool   `json:"value,omitempty"`
	Expires int64  `json:"expires,omitempty"`
	//Context []map[string]string `json:"context,omitempty"`
}

type Group struct {
	Name        string         `json:"name,omitempty"`
	DisplayName string         `json:"displayName,omitempty"`
	Nodes       []*Node        `json:"nodes,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

func (o *Client) GetGroup(ctx context.Context, id string) (*Group, error) {
	resp, err := o.do(ctx, "GET", fmt.Sprintf("group/%s", id), nil)
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
	resp, err := o.do(ctx, "GET", "group", nil)
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

func (o *Client) AddUserToGroup(ctx context.Context, userID string, groupID string, expires *time.Time) (*User, error) {
	node := Node{
		Key:   fmt.Sprintf("group.%s", groupID),
		Value: true,
	}
	if expires != nil {
		node.Expires = expires.Unix()
	}

	body, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}

	_, err = o.do(ctx, "POST", fmt.Sprintf("user/%s/nodes?temporaryNodeMergeStrategy=replace_existing_if_duration_longer", userID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return o.GetUser(ctx, userID)
}

func (o *Client) RemoveUserFromGroup(ctx context.Context, userID string, groupID string, expires *time.Time) (*User, error) {
	node := Node{
		Key:   fmt.Sprintf("group.%s", groupID),
		Value: false,
	}
	if expires != nil {
		node.Expires = expires.Unix()
	}

	body, err := json.Marshal([]*Node{&node})
	if err != nil {
		return nil, err
	}

	_, err = o.do(ctx, "DELETE", fmt.Sprintf("user/%s/nodes", userID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return o.GetUser(ctx, userID)
}
