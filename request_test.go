package jac

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/darrae/jac/internal/jachttp"
)

func Test_message_MarshalRequest(t *testing.T) {
	type fields struct {
		URI     string
		Body    []byte
		Header  http.Header
		Method  Method
		Context context.Context
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			fields: fields{
				URI:     "/test/path",
				Header:  nil,
				Method:  GET,
				Context: context.Background(),
			},
			args: args{
				addr: "https://example.com",
			},
			wantErr: false,
		},
		{
			fields: fields{
				URI:     "/test/path",
				Header:  nil,
				Method:  GET,
				Context: context.Background(),
			},
			args: args{
				addr: "://example.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &message{
				URI:     tt.fields.URI,
				Body:    tt.fields.Body,
				Header:  tt.fields.Header,
				Method:  tt.fields.Method,
				Context: tt.fields.Context,
			}
			got, err := m.MarshalRequest(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("message.MarshalRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			head := http.Header{}
			head.Add("Auth", "Bearer 123")
			jachttp.SetHeaders(got, head)
			if head := got.Header.Get("Auth"); head == "" {
				t.Errorf("message.MarshalRequest().Header = %v", got.Header)
			}
		})
	}
}

func TestGetRequest_Method(t *testing.T) {
	tests := []struct {
		name string
		g    *GetRequest
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GetRequest{}
			if got := g.Method(); got != tt.want {
				t.Errorf("GetRequest.Method() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequest_Body(t *testing.T) {
	tests := []struct {
		name string
		g    *GetRequest
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GetRequest{}
			if got := g.Body(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRequest.Body() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteRequest_Method(t *testing.T) {
	tests := []struct {
		name string
		g    *DeleteRequest
		want string
	}{
		{
			g:    &DeleteRequest{},
			want: "DELETE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &DeleteRequest{}
			if got := g.Method(); got != tt.want {
				t.Errorf("DeleteRequest.Method() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteRequest_Body(t *testing.T) {
	tests := []struct {
		name string
		g    *DeleteRequest
		want []byte
	}{
		{g: &DeleteRequest{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &DeleteRequest{}
			if got := g.Body(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRequest.Body() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteRequest_Header(t *testing.T) {
	tests := []struct {
		name string
		g    *DeleteRequest
		want http.Header
	}{
		{g: &DeleteRequest{}, want: http.Header{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &DeleteRequest{}
			if got := g.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRequest.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutRequest_Method(t *testing.T) {
	tests := []struct {
		name string
		p    *PutRequest
		want string
	}{
		{p: &PutRequest{}, want: "PUT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PutRequest{}
			if got := p.Method(); got != tt.want {
				t.Errorf("PutRequest.Method() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatchRequest_Method(t *testing.T) {
	tests := []struct {
		name string
		p    *PatchRequest
		want string
	}{
		{p: &PatchRequest{}, want: "PATCH"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PatchRequest{}
			if got := p.Method(); got != tt.want {
				t.Errorf("PatchRequest.Method() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutRequest_Header(t *testing.T) {
	tests := []struct {
		name string
		p    *PutRequest
		want http.Header
	}{
		{p: &PutRequest{}, want: http.Header{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PutRequest{}
			if got := p.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutRequest.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatchRequest_Query(t *testing.T) {
	tests := []struct {
		name string
		p    *PatchRequest
		want url.Values
	}{
		{p: &PatchRequest{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PatchRequest{}
			if got := p.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchRequest.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatchRequest_Header(t *testing.T) {
	tests := []struct {
		name string
		p    *PatchRequest
		want http.Header
	}{
		{p: &PatchRequest{}, want: http.Header{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PatchRequest{}
			if got := p.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchRequest.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostRequest_Query(t *testing.T) {
	tests := []struct {
		name string
		p    *PostRequest
		want url.Values
	}{
		{p: &PostRequest{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostRequest{}
			if got := p.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostRequest.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostRequest_Header(t *testing.T) {
	tests := []struct {
		name string
		p    *PostRequest
		want http.Header
	}{
		{p: &PostRequest{}, want: http.Header{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostRequest{}
			if got := p.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostRequest.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutRequest_Query(t *testing.T) {
	tests := []struct {
		name string
		p    *PutRequest
		want url.Values
	}{
		{p: &PutRequest{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PutRequest{}
			if got := p.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutRequest.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequest_Query(t *testing.T) {
	tests := []struct {
		name string
		g    *GetRequest
		want url.Values
	}{
		{g: &GetRequest{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GetRequest{}
			if got := g.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRequest.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteRequest_Query(t *testing.T) {
	tests := []struct {
		name string
		g    *DeleteRequest
		want url.Values
	}{
		{g: &DeleteRequest{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &DeleteRequest{}
			if got := g.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRequest.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
