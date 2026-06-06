package client

import (
	"context"
	"fmt"
	"net/http"
)

type ServerGroup struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Color   string  `json:"color"`
	Servers []int64 `json:"servers"`
}

type CreateServerGroupRequest struct {
	Name    string  `json:"name"`
	Color   string  `json:"color"`
	Servers []int64 `json:"servers"`
}

type UpdateServerGroupRequest struct {
	Name    string  `json:"name,omitempty"`
	Color   string  `json:"color,omitempty"`
	Servers []int64 `json:"servers"`
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

func (c *Client) UpdateServerGroup(ctx context.Context, id int64, req UpdateServerGroupRequest) (*ServerGroup, error) {
	var out ServerGroup

	if err := c.do(ctx, http.MethodPut, fmt.Sprintf("/server/group/%d", id), req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// DeleteServerGroup deletes a server group. The API returns 204 No Content on success.
func (c *Client) DeleteServerGroup(ctx context.Context, id int64) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/server/group/%d", id), nil, nil)
}
