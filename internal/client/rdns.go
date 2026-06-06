package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Rdns struct {
	ID      int64  `json:"id"`
	Version int    `json:"version"`
	IP      string `json:"ip"`
	Rdata   string `json:"rdata"`
	Type    string `json:"type"`
}

type CreateRdnsRequest struct {
	IP    string `json:"ip"`
	Type  string `json:"type"`
	Rdata string `json:"rdata"`
}

type UpdateRdnsRequest struct {
	Rdata string `json:"rdata"`
}

func (c *Client) CreateRdns(ctx context.Context, req CreateRdnsRequest) (*Rdns, error) {
	var out struct {
		Rdns Rdns `json:"rdns"`
	}

	if err := c.do(ctx, http.MethodPost, "/network/rdns", req, &out); err != nil {
		return nil, err
	}

	return &out.Rdns, nil
}

func (c *Client) ListRdns(ctx context.Context, networkID string) ([]Rdns, error) {
	var out []Rdns

	path := "/network/rdns?id=" + url.QueryEscape(networkID)

	if err := c.do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) ReadRdns(ctx context.Context, networkID string, rdnsID int64) (*Rdns, error) {
	entries, err := c.ListRdns(ctx, networkID)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.ID == rdnsID {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("rdns entry %d not found in network %s", rdnsID, networkID)
}

func (c *Client) UpdateRdns(ctx context.Context, id int64, req UpdateRdnsRequest) (*Rdns, error) {
	var out struct {
		Rdns Rdns `json:"rdns"`
	}

	if err := c.do(ctx, http.MethodPut, fmt.Sprintf("/network/rdns/%d", id), req, &out); err != nil {
		return nil, err
	}

	return &out.Rdns, nil
}

func (c *Client) DeleteRdns(ctx context.Context, id int64) error {
	var out struct {
		Rdns Rdns `json:"rdns"`
	}

	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/network/rdns/%d", id), nil, &out)
}
