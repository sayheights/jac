package jac

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/darrae/jac/internal/jachttp"
	"golang.org/x/time/rate"
)

const (
	// defMaxRetryAmount is default retry count used by the client.
	defMaxRetryAmount = 5
)

// Authorizer is the interface implemented by types that can authenticate
// an HTTP request.
type Authorizer interface {
	Authorize(r *http.Request) error
}

// Retry holds the information regarding the retry logic of a Client.
type Retry struct {
	Policy    RetryPolicy
	Backoff   Backoff
	MaxAmount int
}

// Client is an API client capable of making requests to an HTTP based web API.
//
// Client instances have their own retry and authentication logic.
type Client struct {
	Name string
	// BaseURL is the url appended to the path of outgoing requests.
	BaseURL string

	// Headers are the default header values added to every outgoing request.
	Headers http.Header

	// Authorizer handles the authorization of HTTP Request sent from this client.
	Authorizer Authorizer

	// DisableLogging flag determines wheter the client should output logs or not.
	// Default is false.
	DisableLogging bool

	// Retry policy used by the client.
	Retry *Retry

	// Limiter specifies the rate limit.
	Limiter *rate.Limiter

	// TLSConfig is the custom TLSConfig the client will use.
	//
	// If the field is set, the Client will generate a new transport with
	// provided configuration and build a new client with the transport and use that
	// instead of the default client in jachttp package.
	TLSConfig *tls.Config

	Cache Cache

	// IsSuccessful determines if a request should be considered successful or not.
	IsSuccessful func(*http.Response) bool

	hc     *http.Client
	logger txnLogger

	once sync.Once

	mutexes sync.Map
}

func defaultIsSuccessful(res *http.Response) bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}

// init is called once per Client instance and executes the initialization logic.
func (c *Client) init() {
	c.validate()
	c.initRetry()
	c.initLimiter()
	c.initAuth()
	c.initLogger()
	c.initHTTPClient()
	c.initCache()
	if c.IsSuccessful == nil {
		c.IsSuccessful = defaultIsSuccessful
	}
}

func (c *Client) validate() {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		panic("jac: invalid base url: " + c.BaseURL)
	}

	c.BaseURL = strings.TrimSuffix(c.BaseURL, "/")

	if host := baseURL.Host; host == "" {
		panic("jac: base url missing host: " + c.BaseURL)
	}
}

func (c *Client) initRetry() {
	if c.Retry == nil {
		c.Retry = &Retry{
			Policy:    DefaultPolicy,
			Backoff:   DefaultBackoff,
			MaxAmount: defMaxRetryAmount,
		}
	}
}

func (c *Client) initAuth() {
	if c.Authorizer == nil {
		c.Authorizer = zeroAuth
	}
}

func (c *Client) initLogger() {
	if c.logger == nil {
		if c.DisableLogging {
			c.logger = &noopLogger{}

			return
		}
		c.logger = newLogger(c.BaseURL, c.Name)
	}
}

func (c *Client) initHTTPClient() {
	if c.TLSConfig != nil {
		c.hc = jachttp.NewClientWithTLS(c.TLSConfig)
	}
}

// Do makes an HTTP Request with the given context and returns the Response.
func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	c.once.Do(c.init)
	method, ok := strToMethod[req.Method()]
	if !ok {
		return nil, fmt.Errorf("jac: invalid HTTP method %s", req.Method())
	}

	uri := BuildURI(req.Path(), req.Query())
	body := req.Body()
	header := req.Header()
	if header == nil {
		header = http.Header{}
	}

	message := &message{URI: uri, Body: body, Header: header, Method: method, Context: ctx}

	httpRequest, err := c.createRequest(message)
	if err != nil {
		return nil, err
	}

	switch x := req.(type) {
	case CacheRequest:
		return c.doCache(x.CacheKey(), x.TTL(), httpRequest)
	}

	return c.do(httpRequest)
}

func (c *Client) Get(ctx context.Context, u string) (*Response, error) {
	c.once.Do(c.init)
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+u, bytes.NewReader(nil))
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) doCache(key string, dur time.Duration, req *http.Request) (*Response, error) {
	// Wait for inflight requests with the same key to complete
	mutexVal, ok := c.mutexes.LoadOrStore(key, &sync.Mutex{})
	mutex := mutexVal.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	if response := c.Cache.Get(key); response != nil {
		return response, nil
	}

	response, err := c.do(req)
	if err != nil {
		return nil, err
	}

	c.Cache.Set(key, newCacheItem(response, dur))
	if ok {
		c.mutexes.Delete(key)
	}

	return response, nil
}

// createRequest returns an HTTP Request with the given message.
func (c *Client) createRequest(r *message) (*http.Request, error) {
	req, err := r.MarshalRequest(c.BaseURL)
	if err != nil {
		return nil, err
	}

	err = c.prepRequest(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) do(request *http.Request) (*Response, error) {
	err := c.Limiter.Wait(request.Context())
	if err != nil {
		return nil, err
	}
	t, err := c.beginTxn(request)
	if err != nil {
		return nil, err
	}

	c.logger.Log(t)

	state := txnInitial
	for t.reset(); !state.isDone(); {
		state = t.attempt()
		t.reset()
		c.logger.Log(t)
	}

	t.end = time.Now()

	res, err := t.result()
	if err == nil {
		c.logger.Log(t)
	}

	return res, err
}

// prepRequest sets the host, default header values for an HTTP Request
// and authorizes it with the Client's Authorizer.
func (c *Client) prepRequest(req *http.Request) error {
	jachttp.SetHeaders(req, c.Headers)

	return c.Authorizer.Authorize(req)
}

func (c *Client) beginTxn(req *http.Request) (*transaction, error) {
	t := &transaction{req: req, ret: c.Retry, isSuccessful: c.IsSuccessful}
	if c.hc != nil {
		t.hc = c.hc
	}

	return t, t.init()
}

func (c *Client) initCache() {
	if c.Cache == nil {
		c.Cache = NewInMemoryCache()
	}
}

func (c *Client) initLimiter() {
	if c.Limiter == nil {
		c.Limiter = rate.NewLimiter(rate.Inf, 0)
	}
}

func BuildURI(rawPath string, query url.Values) string {
	path := strings.TrimSpace(rawPath)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	path = "/" + path
	removeUnsetParams(query)
	return strings.TrimSuffix(path+"?"+query.Encode(), "?")
}

// removeUnsetParams deletes the keys of the query parameters
// whose values are set to empty string.
func removeUnsetParams(query url.Values) {
	for k, vals := range query {
		if len(vals) == 1 && vals[0] == "" {
			query.Del(k)
		}
	}
}

// noopAuth is a no op implementation of the Authorizer interface.
// It is the default Authorizer implementation used by the client when
// no Authorizer is provided
type noopAuth struct{}

// Authorize implements the Authorizer interface.
func (n *noopAuth) Authorize(r *http.Request) error {
	return nil
}

// zeroAuth is the nil noopAuth instance shared by clients.
var zeroAuth *noopAuth
