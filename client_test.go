package jac

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
)

var TestPath = "/request/path"

var TestMethod = "GET"

type testGetRequest struct {
	*GetRequest
	id  string
	max int
}

func (t *testGetRequest) Path() string {
	return TestPath
}

type testCacheGetRequest struct {
	CacheRequest
}

func (g *testCacheGetRequest) Method() string {
	return TestMethod
}

func (g *testCacheGetRequest) Body() []byte {
	return nil
}

func (g *testCacheGetRequest) Header() http.Header {
	return http.Header{}
}

func (g *testCacheGetRequest) Query() url.Values {
	return nil
}

func (t *testCacheGetRequest) Path() string {
	return TestPath
}

func (t *testCacheGetRequest) TTL() time.Duration {
	return 1 * time.Hour
}

func (t *testCacheGetRequest) CacheKey() string {
	return "cacheKey1"
}

type testCacheRequestSecond struct {
	testCacheGetRequest
}

func (t *testCacheRequestSecond) CacheKey() string {
	return "cacheKey2"
}

func (t *testGetRequest) Query() url.Values {
	u := url.Values{}
	u.Add("query", "param")
	u.Add("id", t.id)
	u.Add("max", strconv.Itoa(t.max))

	return u
}

type testPostRequest struct {
	*PostRequest
	id  string
	max int
}

func (t *testPostRequest) Body() []byte {
	return []byte(`{"response":"success"}`)
}

func (t *testPostRequest) Path() string {
	return TestPath
}

func (t *testPostRequest) Query() url.Values {
	u := url.Values{}
	u.Add("query", "param")
	u.Add("id", t.id)
	u.Add("max", strconv.Itoa(t.max))

	return u
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		req        Request
		failAmount int
		failCode   int
		maxRetry   int
		data       string
		want       *Response
		wantErr    bool
	}{
		{
			ctx:      context.Background(),
			req:      &testGetRequest{id: "successAfterFail", max: 2},
			maxRetry: 5,
			data:     "success",
			want: &Response{
				RequestURI:   "/request/path?id=successAfterFail&max=4&query=param",
				Data:         []byte("success"),
				AttemptCount: 3,
				Duration:     time.Millisecond * 3,
			},
		},
		{
			ctx:      context.Background(),
			req:      &testGetRequest{id: "fail", max: 3},
			maxRetry: 1,
			failCode: 429,
			wantErr:  true,
			want: &Response{
				Data: []byte(""),
			},
		},
		{
			ctx:      context.Background(),
			req:      &testPostRequest{id: "postRequest", max: 2},
			maxRetry: 5,
			failCode: 429,
			want: &Response{
				AttemptCount: 3,
				Duration:     time.Millisecond * 3,
				Data:         []byte(`{"response":"success"}`),
			},
		},
	}

	svr := newTestServer("success")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retry := &Retry{
				Backoff:   LinearBackoff(time.Millisecond * 1),
				Policy:    DefaultPolicy,
				MaxAmount: tt.maxRetry,
			}
			c := &Client{
				BaseURL: svr.URL,
				Retry:   retry,
				Name:    "TestAPI",
			}
			got, err := c.Do(tt.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.want.Equal(got) {
				t.Errorf("Client.Do() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestClient_DoCache(t *testing.T) {
	request := testCacheGetRequest{}
	requestSecond := testCacheRequestSecond{}
	currContext := context.Background()
	svr := newTestServer("success")
	retry := &Retry{
		Backoff:   LinearBackoff(time.Millisecond * 1),
		Policy:    DefaultPolicy,
		MaxAmount: 1,
	}
	c := &Client{
		BaseURL: svr.URL,
		Name:    "TestAPI",
		Retry:   retry,
	}

	responseFirst, err := c.Do(currContext, &request)
	if err != nil {
		t.Errorf("Request execution failed")
		return
	}
	responseSecond, err := c.Do(currContext, &requestSecond)
	if err != nil {
		t.Errorf("Request execution failed")
		return
	}

	tests := []struct {
		name       string
		ctx        context.Context
		req        Request
		failAmount int
		failCode   int
		maxRetry   int
		data       string
		want       *Response
		wantErr    bool
	}{
		{
			name:     "a",
			ctx:      currContext,
			req:      &request,
			maxRetry: 5,
			data:     "success",
			want:     responseFirst,
		},
		{
			name:     "b",
			ctx:      currContext,
			req:      &request,
			maxRetry: 5,
			data:     "success",
			want:     responseFirst,
		},
		{
			name:     "c",
			ctx:      currContext,
			req:      &requestSecond,
			maxRetry: 5,
			data:     "success",
			want:     responseSecond,
		},
		{
			name:     "d",
			ctx:      currContext,
			req:      &requestSecond,
			maxRetry: 5,
			data:     "success",
			want:     responseSecond,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := c.Do(tt.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.want.Equal(got) {
				t.Errorf("Client.Do() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func newTestServer(data string) *httptest.Server {
	handler := &testHandler{
		attempt: map[string]int64{},
		data:    []byte(data),
	}
	svr := httptest.NewServer(handler)
	return svr
}

type testHandler struct {
	attempt map[string]int64
	data    []byte
	mu      sync.Mutex
}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	maxAttempt, _ := strconv.ParseInt(r.URL.Query().Get("max"), 10, 64)
	t.mu.Lock()
	defer t.mu.Unlock()
	attempt := t.attempt[id]

	if attempt < maxAttempt {
		w.WriteHeader(429)
		t.attempt[id] = attempt + 1
		return
	}
	if r.Method == TestMethod || r.Method == "DELETE" {
		fmt.Fprintf(w, string(t.data))
		return
	}
	bod, _ := io.ReadAll(r.Body)
	_, _ = w.Write(bod)
}

func testQuery() url.Values {
	u := url.Values{}
	u.Add("q", "")
	u.Add("c", "a,b,c")

	return u
}

func Test_buildURI(t *testing.T) {
	type args struct {
		rawPath string
		query   url.Values
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				rawPath: "/test/path",
				query:   testQuery(),
			},
			want: "/test/path?c=a%2Cb%2Cc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildURI(tt.args.rawPath, tt.args.query); got != tt.want {
				t.Errorf("buildURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

type failHandler struct {
	last bool
}

func (f *failHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.last {
		w.WriteHeader(404)
		return
	}
	f.last = true
	w.WriteHeader(429)
}

type testLogger struct {
	infoCnt int
	warCnt  int
	errCnt  int
}

func (t *testLogger) Log(txn *transaction) {
	if txn.state == txnInitial {
		t.infoCnt += 1
	}
	if txn.state == txnRetryable {
		t.warCnt += 1
	}
	if txn.state&txnFail != 0 {
		t.errCnt += 1
	}
}

func TestClient_logger(t *testing.T) {
	wantErr := true
	srv := httptest.NewServer(&failHandler{})
	logger := &testLogger{}
	c := &Client{BaseURL: srv.URL, logger: logger}
	c.init()
	defer srv.Close()

	_, err := c.Do(context.Background(), &testGetRequest{})
	if (err != nil) != wantErr {
		t.Errorf("Client.do() error = %v, wantErr %v", err, wantErr)
		return
	}

	if logger.errCnt != 1 || logger.warCnt != 1 || logger.infoCnt != 1 {
		t.Errorf("Client.logger() = %+v", logger)
	}
}

func TestClient_nooplogger(t *testing.T) {
	wantErr := false
	srv := httptest.NewServer(&testHandler{})
	c := &Client{BaseURL: srv.URL, DisableLogging: true}
	c.init()
	defer srv.Close()

	_, err := c.Do(context.Background(), &testGetRequest{})
	if (err != nil) != wantErr {
		t.Errorf("Client.do() error = %v, wantErr %v", err, wantErr)
		return
	}
}
