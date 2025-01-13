package jac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

type testPaginatedRequest struct {
	*GetRequest
	page int
}

func (t *testPaginatedRequest) Path() string {
	return "/"
}

func (t *testPaginatedRequest) Query() url.Values {
	u := url.Values{}
	u.Add("page", strconv.Itoa(t.page))

	return u
}

func (t *testPaginatedRequest) Body() []byte {
	return nil
}

func (t *testPaginatedRequest) Next(p *Response) (PaginatedRequest, bool) {
	var response paginationResponse
	err := json.Unmarshal(p.Data, &response)
	if err != nil {
		return nil, true
	}

	if response.Count == 0 {
		return nil, true
	}
	t.page = t.page + 1

	return t, false
}

type paginationHandler struct{}

type paginationResponse struct {
	Count   int
	Success bool
}

func (p *paginationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := `{"Count":10}`
	done := `{"Count":0,"Success":true}`
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	w.WriteHeader(200)
	if page == 4 {
		fmt.Fprintf(w, done)
		return
	}
	fmt.Fprintf(w, msg)
	return
}

func TestClient_DoPagination(t *testing.T) {
	tests := []struct {
		name    string
		p       PaginatedRequest
		want    int
		wantErr bool
	}{
		{
			p:    &testPaginatedRequest{},
			want: 5,
		},
		{
			p:    &testPaginatedRequest{page: 3},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &paginationHandler{}
			srv := httptest.NewServer(handler)
			c := &Client{BaseURL: srv.URL}
			defer srv.Close()
			got, err := c.DoPagination(context.Background(), tt.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DoPagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Client.DoPagination() count total = %d, want %d", len(got), tt.want)
			}
			for i, res := range got {
				var resp paginationResponse
				_ = json.Unmarshal(res.Data, &resp)
				if i == len(got)-1 {
					if resp.Count != 0 || !resp.Success {
						t.Errorf("Client.DoPagination() count %d, success %t", resp.Count, resp.Success)
						return
					}
					return
				}
				if resp.Count != 10 {
					t.Errorf("Client.DoPagination() count = %d", resp.Count)
					return
				}

			}
		})
	}
}
