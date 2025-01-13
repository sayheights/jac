package auth

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type BearerAPI struct {
	ClientSecret string
	RefreshToken string
	URL          string
	ExpiresIn    int64
	hasAPIKey    bool
	APIKey       string
	tok          *Token
	start        time.Time

	clientSecretKey string
	refreshTokenKey string
	once            sync.Once
	mu              sync.Mutex
}

type ClientSecret struct {
	HeaderKey string
	Value     string
}

type RefreshToken struct {
	HeaderKey string
	Value     string
}

func NewBearerAPI(cs ClientSecret, rt RefreshToken, authURL string) *BearerAPI {
	return &BearerAPI{
		ClientSecret:    cs.Value,
		clientSecretKey: cs.HeaderKey,
		RefreshToken:    rt.Value,
		refreshTokenKey: rt.HeaderKey,
		URL:             authURL,
	}
}

func (c *BearerAPI) hasExpired() bool {
	return c.start.Add(time.Second * time.Duration(c.tok.ExpiresIn)).Before(time.Now())
}

// BuildRequest creates the request which returns access token
func (c *BearerAPI) BuildRequest() (*http.Request, error) {
	c.once.Do(func() {
		if c.clientSecretKey == "" {
			c.clientSecretKey = "secretkey"
		}
		if c.refreshTokenKey == "" {
			c.refreshTokenKey = "refreshToken"
		}
	})
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(c.clientSecretKey, c.ClientSecret)
	req.Header.Add(c.refreshTokenKey, c.RefreshToken)
	return req, nil
}

// Authorize Handles bearer api cases
// There are two cases in bearer api
// If it has api key it simply add header api key with bearer token type
// If not it builds a request and execution of request provides us with access token and it adds access token to header
// If token is already exist and not expired we simply use existing token
func (c *BearerAPI) Authorize(r *http.Request) error {
	if c.tok == nil || c.hasExpired() {
		c.mu.Lock()
		defer c.mu.Unlock()
		var t Token
		if c.hasAPIKey {
			t = Token{
				AccessToken:  strings.TrimPrefix(strings.TrimSuffix(c.APIKey, "\""), "\""),
				ExpiresIn:    int(c.ExpiresIn),
				RefreshToken: c.RefreshToken,
				TokenType:    "Bearer",
			}
		} else {
			c.start = time.Now()
			req, err := c.BuildRequest()
			if err != nil {
				return err
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			accessToken, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			t = Token{
				AccessToken:  strings.TrimPrefix(strings.TrimSuffix(string(accessToken), "\""), "\""),
				ExpiresIn:    int(c.ExpiresIn),
				RefreshToken: c.RefreshToken,
				TokenType:    "Bearer",
			}
		}
		c.tok = &t
		r.Header.Add("Authorization", c.tok.format())
		return nil
	}
	r.Header.Add("Authorization", c.tok.format())
	return nil
}
