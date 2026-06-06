package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type Network struct {
	ID           string `json:"id"`
	CIDR         string `json:"cidr"`
	Version      int    `json:"version"`
	Gateway      string `json:"gateway"`
	Network      string `json:"network"`
	Broadcast    string `json:"broadcast"`
	Netmask      string `json:"netmask"`
	PrefixLength int    `json:"prefixLength"`
	Registry     string `json:"registry"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Resolver     string `json:"resolver"`
}

type ListNetworksFilter struct {
	Version  int
	Registry string
	Country  string
	CIDR     string
	IP       string
}

func (c *Client) ListNetworks(ctx context.Context, f ListNetworksFilter) ([]Network, error) {
	var out []Network

	q := url.Values{}
	if f.Version != 0 {
		q.Set("IP version", strconv.Itoa(f.Version))
	}
	if f.Registry != "" {
		q.Set("registry", f.Registry)
	}
	if f.Country != "" {
		q.Set("country", f.Country)
	}
	if f.CIDR != "" {
		q.Set("cidr", f.CIDR)
	}
	if f.IP != "" {
		q.Set("ip", f.IP)
	}

	path := "/network"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	if err := c.do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}

	return out, nil
}

