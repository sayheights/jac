package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestOAuth2_Authorize(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock()
	// Requests we expect to change with Authorize method
	// token we expect to not change with Authorize method
	bodyHeader := make(map[string][]string)
	bodyRequest := &http.Request{
		Method: "GET",
		Header: bodyHeader,
	}
	formHeader := make(map[string][]string)
	formRequest := &http.Request{
		Method: "GET",
		Header: formHeader,
	}
	tok := Token{
		AccessToken:  "Bearer existingToken",
		TokenType:    "",
		ExpiresIn:    int(time.Now().Unix() + 100),
		RefreshToken: "",
	}
	type fields struct {
		ClientId     string
		ClientSecret string
		RefreshToken string
		GrantType    string
		URL          string
		IsBody       bool
		IsForm       bool
		IsHeader     bool
		ExpiresIn    int64
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
			name: "body",
			fields: fields{
				ClientId:     "bodyId",
				ClientSecret: "bodySecret",
				RefreshToken: "bodyToken",
				GrantType:    "",
				URL:          "https://127.0.0:3000/body",
				IsBody:       true,
				IsForm:       false,
				ExpiresIn:    0,
				tok:          nil,
				start:        time.Time{},
			},
			args: args{
				r: bodyRequest,
			},
			wantErr: false,
		},
		{
			name: "form",
			fields: fields{
				ClientId:     "formId",
				ClientSecret: "formSecret",
				RefreshToken: "formToken",
				GrantType:    "",
				URL:          "https://127.0.0:3000/form",
				IsBody:       false,
				IsForm:       true,
				ExpiresIn:    0,
				tok:          nil,
				start:        time.Time{},
			},
			args: args{
				r: formRequest,
			},
			wantErr: false,
		},
		{
			name: "existing token",
			fields: fields{
				ClientId:     "bodyId",
				ClientSecret: "bodySecret",
				RefreshToken: "bodyToken",
				GrantType:    "",
				URL:          "https://127.0.0:3000/body",
				ExpiresIn:    time.Now().Unix() + 100,
				tok:          &tok,
				start:        time.Time{},
			},
			args: args{
				r: bodyRequest,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &OAuth2{
				ClientId:     tt.fields.ClientId,
				ClientSecret: tt.fields.ClientSecret,
				RefreshToken: tt.fields.RefreshToken,
				GrantType:    tt.fields.GrantType,
				URL:          tt.fields.URL,
				IsBody:       tt.fields.IsBody,
				IsForm:       tt.fields.IsForm,
				ExpiresIn:    tt.fields.ExpiresIn,
				tok:          tt.fields.tok,
				start:        tt.fields.start,
			}
			req := tt.args.r
			if err := c.Authorize(req); (err != nil) != tt.wantErr {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	bodyValue, ok := bodyHeader["Authorization"]
	if !ok {
		t.Errorf("Authorize() method could not create body token")
	}
	formValue, ok := formHeader["Authorization"]
	if !ok {
		t.Errorf("Authorize() method could not create form token")
	}

	assert.Equal(t, bodyValue[0], "Bearer bodyToken", "Authorize() method create body token falsely")
	assert.Equal(t, formValue[0], "Bearer formToken", "Authorize() method create form token falsely")
	assert.Equal(t, tok.AccessToken, "Bearer existingToken", "Authorize() method changes existing token even it is valid")
}

func mock() {
	httpmock.RegisterResponder(
		"POST",
		"https://127.0.0:3000/body",
		httpmock.NewStringResponder(200,
			"{\"access_token\":\"bodyToken\", \"token_type\":\"Bearer\", \"ExpiresIn\":0, \"RefreshToken\":\"\"}"),
	)
	httpmock.RegisterResponder(
		"POST",
		"https://127.0.0:3000/form",
		httpmock.NewStringResponder(200,
			"{\"access_token\":\"formToken\", \"token_type\":\"Bearer\", \"ExpiresIn\":0, \"RefreshToken\":\"\"}"),
	)
}
