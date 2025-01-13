package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic_Authorize(t *testing.T) {
	// Request we expect to change with Authorize method
	req, _ := http.NewRequest("GET", "https://127.0.0:3000", nil)
	assert.Empty(t, req.Header)
	type fields struct {
		ID  string
		Key string
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
			name: "basic authentication",
			fields: fields{
				ID:  "basic",
				Key: "basicKey",
			},
			args: args{
				req,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Basic{
				ID:  tt.fields.ID,
				Key: tt.fields.Key,
			}
			if err := b.Authorize(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	assert.NotEmpty(t, req.Header)
}
