package jac

import (
	"bytes"
	"net/http"
	"net/url"
	"time"
)

// Response represents an HTTP Response.
type Response struct {
	URL          *url.URL
	RequestURI   string
	Data         []byte
	Header       http.Header
	Duration     time.Duration
	AttemptCount int
	StatusCode   int
}

func (r *Response) Equal(o *Response) bool {
	if o == nil {
		return r == nil
	}
	att := r.AttemptCount == o.AttemptCount
	dur := r.Duration.Round(time.Second) == o.Duration.Round(time.Second)
	if att && dur {
		return bytes.Equal(r.Data, o.Data)
	}
	return false
}
