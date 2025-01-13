package auth

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestBearerAPI_Authorize(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockBearer()

	// Requests we expect to change with Authorize method
	apiKeyHeader := make(map[string][]string)
	apiKeyRequest := &http.Request{
		Method: "GET",
		Header: apiKeyHeader,
	}
	bearerHeader := make(map[string][]string)
	bearerRequest := &http.Request{
		Method: "GET",
		Header: bearerHeader,
	}

	// token we expect to not change with Authorize method
	tok := Token{
		AccessToken:  "Bearer existingToken",
		TokenType:    "",
		ExpiresIn:    int(time.Now().Unix() + 100),
		RefreshToken: "",
	}
	type fields struct {
		ClientSecret string
		RefreshToken string
		URL          string
		ExpiresIn    int64
		hasAPIKey    bool
		APIKey       string
		tok          *Token
		start        time.Time
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
			name: "api key",
			fields: fields{
				ClientSecret: "apiKeySecret",
				RefreshToken: "apiKeyToken",
				ExpiresIn:    0,
				tok:          nil,
				start:        time.Time{},
				hasAPIKey:    true,
				APIKey:       "apiKeyToken",
			},
			args: args{
				r: apiKeyRequest,
			},
			wantErr: false,
		},
		{
			name: "bearer",
			fields: fields{
				ClientSecret: "bearerSecret",
				RefreshToken: "bearerToken",
				URL:          "https://127.0.0:3000/bearer",
				ExpiresIn:    0,
				tok:          nil,
				start:        time.Time{},
				hasAPIKey:    false,
			},
			args: args{
				r: bearerRequest,
			},
			wantErr: false,
		},
		{
			name: "false url",
			fields: fields{
				ClientSecret: "falseSecret",
				RefreshToken: "bearerToken",
				URL:          "https://127.0.0:3000/falsebearer",
				ExpiresIn:    0,
				tok:          nil,
				start:        time.Time{},
				hasAPIKey:    false,
			},
			args: args{
				r: bearerRequest,
			},
			wantErr: true,
		},
		{
			name: "error response",
			fields: fields{
				ClientSecret: "errorSecret",
				RefreshToken: "bearerToken",
				URL:          "https://127.0.0:3000/error",
				ExpiresIn:    0,
				tok:          nil,
				start:        time.Time{},
				hasAPIKey:    false,
			},
			args: args{
				r: bearerRequest,
			},
			wantErr: true,
		},
		{
			name: "existing token",
			fields: fields{
				ClientSecret: "apiKeySecret",
				RefreshToken: "apiKeyToken",
				ExpiresIn:    time.Now().Unix() + 100,
				tok:          &tok,
				start:        time.Time{},
				hasAPIKey:    true,
				APIKey:       "apiKeyToken",
			},
			args: args{
				r: apiKeyRequest,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &BearerAPI{
				ClientSecret: tt.fields.ClientSecret,
				RefreshToken: tt.fields.RefreshToken,
				URL:          tt.fields.URL,
				ExpiresIn:    tt.fields.ExpiresIn,
				hasAPIKey:    tt.fields.hasAPIKey,
				APIKey:       tt.fields.APIKey,
				tok:          tt.fields.tok,
				start:        tt.fields.start,
			}
			if err := c.Authorize(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	bearerValue, ok := bearerHeader["Authorization"]
	if !ok {
		t.Errorf("Authorize() method could not create bearer access token")
	}
	apiKeyValue, ok := apiKeyHeader["Authorization"]
	if !ok {
		t.Errorf("Authorize() method could not create api key token")
	}

	assert.Equal(t, bearerValue[0], "Bearer bearerAccessToken", "Authorize() method create bearer access token falsely")
	assert.Equal(t, apiKeyValue[0], "Bearer apiKeyToken", "Authorize() method create api key token falsely")
	assert.Equal(t, tok.AccessToken, "Bearer existingToken", "Authorize() method changes existing token even it is valid")
}

func mockBearer() {
	httpmock.RegisterResponder(
		"GET",
		"https://127.0.0:3000/bearer",
		httpmock.NewStringResponder(200,
			"bearerAccessToken"),
	)
	origError := errors.New("body error")
	httpmock.RegisterResponder(
		"GET",
		"https://127.0.0:3000/error",
		httpmock.NewErrorResponder(origError),
	)
}
