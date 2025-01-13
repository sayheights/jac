package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIKey_Authorize(t *testing.T) {
	// Requests we expect to change with Authorize method
	req, _ := http.NewRequest("GET", "https://127.0.0:3000", nil)
	reqHeader, _ := http.NewRequest("GET", "https://127.0.0:3000", nil)
	// First check if they initially empty
	assert.Empty(t, req.URL.RawQuery)
	assert.Empty(t, reqHeader.Header)
	type fields struct {
		Key   string
		Value string
		In    In
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
			name: "api key auth in query",
			fields: fields{
				Key:   "queryKey",
				Value: "queryValue",
				In:    InQuery,
			},
			args: args{
				req,
			},
			wantErr: false,
		},
		{
			name: "api key auth in header",
			fields: fields{
				Key:   "headerKey",
				Value: "headerValue",
				In:    InHeader,
			},
			args: args{
				reqHeader,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIKey{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
				In:    tt.fields.In,
			}
			if err := a.Authorize(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	headerValue, ok := reqHeader.Header["Headerkey"]
	if !ok {
		t.Errorf("Authorize() method could not create header value")
	}
	queryValue := req.URL.RawQuery
	if !ok {
		t.Errorf("Authorize() method could not create query value")
	}
	// The changes that should happen
	assert.Equal(t, headerValue[0], "headerValue", "Authorize() method calculate header value falsely")
	assert.Equal(t, queryValue, "queryKey=queryValue", "Authorize() method calculate query value falsely")
}
