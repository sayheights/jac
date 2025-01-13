package jac

import (
	"context"
	"sync"
	"time"
)

type AsyncRequest interface {
	Request
	IsReady(*Response) (bool, error)
	OnReady() Request
}

type AsyncResponse struct {
	Response *Response
	Err      error
}

// DoAsync repeats a given request until the provided Until function returns true.
func (c *Client) DoAsync(ctx context.Context, asyncs ...AsyncRequest) chan *AsyncResponse {
	ch := make(chan *AsyncResponse)
	go c.doAsyncs(ch, ctx, asyncs)

	return ch
}

func (c *Client) doAsyncs(ch chan *AsyncResponse, ctx context.Context, asyncs []AsyncRequest) {
	var wg sync.WaitGroup
	for _, async := range asyncs {
		a := async
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.doAsync(ctx, a, ch)
		}()
	}
	wg.Wait()
	close(ch)
}

func (c *Client) doAsync(ctx context.Context, request AsyncRequest, ch chan *AsyncResponse) {
	var attempt int
	for {
		if attempt > 5 {
			attempt = 0
		}

		if cacheReq, ok := isAsyncReadyCached(request); ok {
			response := c.Cache.Get(cacheReq.CacheKey())
			if response != nil {
				ch <- &AsyncResponse{Response: response}
				return
			}
		}

		response, err := c.Do(ctx, request)
		if err != nil {
			ch <- &AsyncResponse{Err: err}
			return
		}

		ok, err := request.IsReady(response)
		if err != nil {
			ch <- &AsyncResponse{Err: err}
			return
		}

		if ok {
			c.doAsyncReady(ctx, request.OnReady(), ch)
			return
		}
		attempt += 1

		time.Sleep(c.Retry.Backoff(attempt))
	}
}

func (c *Client) doAsyncReady(ctx context.Context, request Request, ch chan *AsyncResponse) {
	res, err := c.Do(ctx, request)
	ch <- &AsyncResponse{Response: res, Err: err}
	return
}

func isAsyncReadyCached(req AsyncRequest) (CacheRequest, bool) {
	switch x := req.(type) {
	case CacheRequest:
		return x, true

	default:
		return nil, false
	}
}
