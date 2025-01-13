package auth

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOAuth_Authorize(t *testing.T) {
	// Request we expect to change with Authorize method
	oauthRequest := &http.Request{
		Method: "GET",
		Header: make(map[string][]string),
		URL:    &url.URL{Host: "https://127.0.0:3000", Path: "oauth"},
	}
	type fields struct {
		ConsumerKey    string
		ConsumerSecret string
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "general test",
			fields: fields{
				ConsumerKey:    "consKey",
				ConsumerSecret: "consSecret",
			},
			args: args{
				r: oauthRequest,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &OAuth{
				ConsumerKey:    tt.fields.ConsumerKey,
				ConsumerSecret: tt.fields.ConsumerSecret,
			}
			if err := c.Authorize(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	assert.NotEmpty(t, oauthRequest.Header)
}
