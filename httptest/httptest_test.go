package httptest

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

type testRequest struct{}

func (t *testRequest) Method() string {
	return "POST"
}

func (t *testRequest) Path() string {
	return "/resource/path"
}

func (t *testRequest) Query() url.Values {
	u := url.Values{}
	u.Add("key", "param")
	u.Add("keys", "param1,param2")

	return u
}

func (t *testRequest) Body() []byte {
	return []byte("REQUEST_BODY")
}

func (t *testRequest) Header() http.Header {
	u := http.Header{}
	u.Add("Content-Type", "text/plain")

	return u
}

func TestServerAndClient(t *testing.T) {
	type testCase struct {
		asserts []Assert
		wantErr bool
	}

	method := "POST"
	path := "/resource/path"
	headers := http.Header{}
	headers.Add("Content-Type", "text/plain")
	headers.Add("Accept-Encoding", "gzip")
	query := "keys=param1,param2&key=param"
	body := []byte("REQUEST_BODY")

	assertMethod := AssertMethod(method)
	assertPath := AssertPath(path)
	assertQuery := AssertQuery(query)
	assertHeaders := AssertHeaders(headers)
	assertBody := AssertBody(body)

	tests := []testCase{
		{asserts: []Assert{assertBody, assertHeaders}},
		{asserts: []Assert{assertMethod, assertQuery}},
		{asserts: []Assert{assertMethod, assertPath, assertQuery, assertHeaders, assertBody}},
		{asserts: []Assert{AssertMethod("GET"), assertQuery}, wantErr: true},
		{asserts: []Assert{assertQuery, AssertMethod("PUT")}, wantErr: true},
		{asserts: []Assert{AssertBody([]byte("REQST_BODY"))}, wantErr: true},
	}

	for i, tt := range tests {
		srv, client := ServerAndClient(SuccessHandler, tt.asserts...)
		_, err := client.Do(context.Background(), &testRequest{})
		if (err != nil) != tt.wantErr {
			t.Errorf("%d: ServerAndClient() error = %v, wantErr %v", i, err, tt.wantErr)
		}
		srv.Close()
	}
}
