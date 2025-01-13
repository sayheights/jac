package auth

import (
	"fmt"
	"strings"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (t *Token) format() string {
	typ := strings.Title(t.TokenType)

	return fmt.Sprintf("%s %s", typ, t.AccessToken)
}
