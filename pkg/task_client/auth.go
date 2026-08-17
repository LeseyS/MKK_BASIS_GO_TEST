package task_client

import (
	"context"
	"net/http"
)

func (c *Client) Register(ctx context.Context, in RegisterRequest) (RegisterResponse, error) {
	var out RegisterResponse

	err := c.do(ctx, http.MethodPost, apiPrefix+"/register", nil, in, &out)

	return out, err
}

func (c *Client) Login(ctx context.Context, in LoginRequest) (LoginResponse, error) {
	var out LoginResponse

	err := c.do(ctx, http.MethodPost, apiPrefix+"/login", nil, in, &out)

	return out, err
}

func (c *Client) RegisterAndLogin(ctx context.Context, name, email, password string) (*Client, User, error) {
	if _, err := c.Register(ctx, RegisterRequest{Name: name, Email: email, Password: password}); err != nil {
		return nil, User{}, err
	}

	out, err := c.Login(ctx, LoginRequest{Email: email, Password: password})
	if err != nil {
		return nil, User{}, err
	}

	return c.WithToken(out.Token), out.User, nil
}
