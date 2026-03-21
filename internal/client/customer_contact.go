package client

///////////////////////////////////////////////////////////////////////////////////MODULES
import (
	"context"
	"fmt"
	"net/http"
)

//////////////////////////////////////////////////////////////////////////////////////////

type CustomerContact struct {
	ID    int64  `json:"Id"`
	Email string `json:"Email"`
	Type  string `json:"Type"`
}

type CreateCustomerContactRequest struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

func (c *Client) CreateCustomerContact(ctx context.Context, req CreateCustomerContactRequest) (*CustomerContact, error) {
	var out CustomerContact

	if err := c.do(ctx, http.MethodPost, "/customer/contact", req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) ReadCustomerContact(ctx context.Context, customerContactID int64) (*CustomerContact, error) {
	var out CustomerContact

	path := fmt.Sprintf("/customer/contact/%d", customerContactID)

	if err := c.do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) DeleteCustomerContact(ctx context.Context, id int64) (*CustomerContact, error) {
	var out CustomerContact

	if err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/customer/contact/%d", id), nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
