package jac

import (
	"context"
	"fmt"
	"strings"
)

type ClientGroup struct {
	clients []*Client
}

func NewClientGroup(clients ...*Client) *ClientGroup {
	return &ClientGroup{clients: clients}
}

func (c *ClientGroup) Get(ctx context.Context, u string) (*Response, error) {
	client := c.getClient(u)
	if client == nil {
		return nil, fmt.Errorf("jac: no matching client found for %s", u)
	}

	return client.Get(context.Background(), u[len(client.BaseURL)+1:])
}

func (c *ClientGroup) getClient(u string) *Client {
	for _, client := range c.clients {
		if strings.HasPrefix(u, client.BaseURL) {
			return client
		}
	}

	return nil
}
