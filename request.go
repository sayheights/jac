package jac

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/darrae/jac/internal/jachttp"
)

type CacheRequest interface {
	Request
	CacheKey() string
	TTL() time.Duration
	EvictionPolicy() func(cache Cache) bool
}

// Request is the interface implemented by types that can be converted into an HTTP Request.
type Request interface {
	Method() string
	Path() string
	Query() url.Values
	Body() []byte
	Header() http.Header
}

// GetRequest is a convenience struct for types that represent an HTTP request with Get Method.
//
// Types that implement Request interface with GET
// method can embed a nil instance of GetRequest.
type GetRequest struct{}

// Method returns the GET http method.
func (g *GetRequest) Method() string {
	return "GET"
}

// Body returns a nil slice as GET requests do not have a body.
func (g *GetRequest) Body() []byte {
	return nil
}

func (g *GetRequest) Header() http.Header {
	return http.Header{}
}

func (g *GetRequest) Query() url.Values {
	return nil
}

type DeleteRequest struct{}

// Method returns the GET http method.
func (g *DeleteRequest) Method() string {
	return "DELETE"
}

// Body returns a nil slice as GET requests do not have a body.
func (g *DeleteRequest) Body() []byte {
	return nil
}

func (g *DeleteRequest) Query() url.Values {
	return nil
}

func (g *DeleteRequest) Header() http.Header {
	return http.Header{}
}

type PostRequest struct{}

func (p *PostRequest) Method() string {
	return "POST"
}

type PutRequest struct{}

func (p *PutRequest) Method() string {
	return "PUT"
}

type PatchRequest struct{}

func (p *PatchRequest) Method() string {
	return "PATCH"
}

func (p *PutRequest) Header() http.Header {
	return http.Header{}
}

func (p *PatchRequest) Query() url.Values {
	return nil
}

func (p *PatchRequest) Header() http.Header {
	return http.Header{}
}

func (p *PostRequest) Query() url.Values {
	return nil
}

func (p *PostRequest) Header() http.Header {
	return http.Header{}
}

func (p *PutRequest) Query() url.Values {
	return nil
}

// message represents an HTTP message. It implements the RequestMarshaler interface.
type message struct {
	URI     string
	Body    []byte
	Header  http.Header
	Method  Method
	Context context.Context
}

// MarsahlRequest marshals the message instance into an HTTP Request
// targeting the given address.
func (m *message) MarshalRequest(addr string) (*http.Request, error) {
	method := m.Method.String()
	url := addr + m.URI
	body := bytes.NewReader(m.Body)

	req, err := http.NewRequestWithContext(m.Context, method, url, body)
	if err != nil {
		return nil, err
	}

	jachttp.SetHeaders(req, m.Header)

	return req, nil
}
