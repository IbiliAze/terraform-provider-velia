package client

import (
	"context"
	"fmt"
	"net/http"
)

type Ticket struct {
	ID         int64    `json:"id"`
	Queue      string   `json:"queue"`
	Status     string   `json:"status"`
	Subject    string   `json:"subject"`
	Priority   int      `json:"priority"`
	Requestors []string `json:"requestors"`
	CC         []string `json:"cc"`
	Created    string   `json:"created"`
	Due        string   `json:"due"`
	Resolved   string   `json:"resolved"`
	Updated    string   `json:"updated"`
	Servers    []int64  `json:"server"`
}

type CreateTicketRequest struct {
	Topic   string  `json:"topic"`
	Subject string  `json:"subject"`
	Message string  `json:"message"`
	Servers []int64 `json:"server,omitempty"`
}

func (c *Client) CreateTicket(ctx context.Context, req CreateTicketRequest) (*Ticket, error) {
	var out Ticket

	if err := c.do(ctx, http.MethodPost, "/ticket", req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) ReadTicket(ctx context.Context, id int64) (*Ticket, error) {
	var out struct {
		Ticket Ticket `json:"ticket"`
	}

	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/ticket/%d", id), nil, &out); err != nil {
		return nil, err
	}

	return &out.Ticket, nil
}
