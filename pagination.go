package jac

import (
	"context"
)

// PaginatedRequest is the interface implemented by Requests that support pagination.
//
// PaginatedRequest embeds a regular Request which holds the Request generation behaviour.
// Next method represents the capability where give a Response a Request can generate the
// Request which would retrieve the next page of data.
type PaginatedRequest interface {
	Request

	// Next accepts a Response as an input and returns the Request
	// for accessing the next page of pagination.
	// Along with the PaginatedRequest Next also returns a boolean
	// flag representing wheter the pagination requests are done.
	// If there are no more pages left, Next returns true.
	Next(*Response) (PaginatedRequest, bool)
}

func (c *Client) DoPagination(ctx context.Context, p PaginatedRequest) ([]*Response, error) {
	var responses []*Response
	var done bool
	for !done {
		res, err := c.Do(ctx, p)
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
		p, done = p.Next(res)
	}

	return responses, nil
}
