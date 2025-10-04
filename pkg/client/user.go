package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type User struct {
	Username string  `json:"username"`
	UniqueID string  `json:"uniqueId"`
	Nodes    []*Node `json:"nodes"`
}

func (o *Client) GetUser(ctx context.Context, id string) (*User, error) {
	resp, err := o.do(ctx, "GET", fmt.Sprintf("user/%s", id), nil)
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

	u := User{}
	if err := json.Unmarshal(body, &u); err != nil {
		if err := resp.Body.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (o *Client) ListAllUsers(ctx context.Context) ([]*User, error) {
	resp, err := o.do(ctx, "GET", "user", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	userIDs := []string{}
	if err := json.Unmarshal(body, &userIDs); err != nil {
		return nil, err
	}

	users := []*User{}
	for _, user := range userIDs {
		u, err := o.GetUser(ctx, user)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (o *Client) ListAllUsersInGroup(ctx context.Context, groupID string) ([]*User, error) {
	resp, err := o.do(ctx, "GET", fmt.Sprintf("user/search?group=%s", groupID), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	users := []*User{}
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}

	return users, nil
}
