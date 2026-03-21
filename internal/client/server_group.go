package client

import (
	"context"
	"fmt"
	"net/http"
)

type ServerGroup struct {
	ID      int64    `json:"Id"`
	Name    string   `json:"Name"`
	Color   string   `json:"Color"`
	Servers []string `json:"Servers"`
}

type CreateServerGroupRequest struct {
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Servers []string `json:"servers"`
}

func (c *Client) CreateServerGroup(ctx context.Context, req CreateServerGroupRequest) (*ServerGroup, error) {
	var out ServerGroup

	if err := c.do(ctx, http.MethodPost, "/server/group", req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) ReadServerGroup(ctx context.Context, serverGroupId int64) (*ServerGroup, error) {
	var out ServerGroup

	path := fmt.Sprintf("/server/group/%d", serverGroupId)

	if err := c.do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) DeleteServerGroup(ctx context.Context, id int64) (*ServerGroup, error) {
	var out ServerGroup

	if err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/server/group/%d", id), nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
