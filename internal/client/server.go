package client

import (
	"context"
	"fmt"
	"net/http"
)

type ServerNic struct {
	ID                    string  `json:"id"`
	Online                bool    `json:"online"`
	TransferRate          int64   `json:"transferRate"`
	TrafficTotalMonthlyTB float64 `json:"trafficTotalMonthlyTB"`
}

type Server struct {
	ID         int64       `json:"id"`
	Label      string      `json:"label"`
	Ticket     int64       `json:"ticket"`
	Status     string      `json:"status"`
	Issued     string      `json:"issued"`
	Cancelled  string      `json:"cancelled"`
	City       string      `json:"city"`
	Country    string      `json:"country"`
	IPMiIP     string      `json:"ipmiIp"`
	Online     bool        `json:"online"`
	Product    string      `json:"product"`
	Type       string      `json:"type"`
	DeviceType string      `json:"deviceType"`
	Tags       []string    `json:"tags"`
	Networks   []string    `json:"networks"`
	ServerIP   []string    `json:"serverIp"`
	Nics       []ServerNic `json:"nics"`
}

type UpdateServerLabelRequest struct {
	Label string `json:"label"`
}

func (c *Client) ReadServer(ctx context.Context, id int64) (*Server, error) {
	var out Server

	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/server/%d", id), nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) UpdateServerLabel(ctx context.Context, id int64, label string) (*Server, error) {
	var out Server

	req := UpdateServerLabelRequest{Label: label}

	if err := c.do(ctx, http.MethodPut, fmt.Sprintf("/server/%d/label", id), req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
