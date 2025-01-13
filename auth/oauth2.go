package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type OAuth2 struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
	GrantType    string
	URL          string
	IsBody       bool
	IsForm       bool
	ExpiresIn    int64
	tok          *Token
	start        time.Time
	mu           sync.Mutex
}

type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	GrantType    string
	IsBody       bool
}

func NewOAuth2(config OAuth2Config, authURL string) *OAuth2 {
	return &OAuth2{
		ClientId:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RefreshToken: config.RefreshToken,
		GrantType:    config.GrantType,
		URL:          authURL,
		IsBody:       config.IsBody,
	}
}

func (c *OAuth2) hasExpired() bool {
	return c.start.Add(time.Second * time.Duration(c.tok.ExpiresIn)).Before(time.Now())
}

// BuildRequest creates the request which returns access token
// There are three cases
// Api can need client credentials in body, form or url
func (c *OAuth2) BuildRequest() (*http.Request, error) {
	if c.IsBody {
		body := map[string]string{
			"client_id":     c.ClientId,
			"client_secret": c.ClientSecret,
			"refresh_token": c.RefreshToken,
			"grant_type":    c.GrantType,
		}

		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		return http.NewRequest("POST", c.URL, bytes.NewReader(b))
	}
	v := url.Values{
		"client_id":     {c.ClientId},
		"client_secret": {c.ClientSecret},
		"refresh_token": {c.RefreshToken},
		"grant_type":    {c.GrantType},
	}
	if c.IsForm {
		req, err := http.NewRequest("POST", c.URL, strings.NewReader(v.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		return req, nil
	}

	c.URL += "?" + v.Encode()
	return http.NewRequest("POST", c.URL, nil)
}

// Authorize Handles oauth2 cases
// It sets access token to header
// If token is already exist and not expired we simply use existing token
func (c *OAuth2) Authorize(r *http.Request) error {
	if c.tok == nil || c.hasExpired() {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.start = time.Now()
		req, err := c.BuildRequest()
		if err != nil {
			return err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		var t Token
		err = json.NewDecoder(res.Body).Decode(&t)
		if err != nil {
			return err
		}
		c.tok = &t
		r.Header.Set("Authorization", c.tok.format())
		return nil
	}
	r.Header.Add("Authorization", c.tok.format())
	return nil
}
