package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OAuth struct {
	ConsumerKey    string
	ConsumerSecret string
}

func NewOAuth(consumerKey, consumerSecret string) *OAuth {
	return &OAuth{ConsumerKey: consumerKey, ConsumerSecret: consumerSecret}
}

// Authorize handles oauth case
// Oauth needs client credentials and signature in header
func (c *OAuth) Authorize(r *http.Request) error {
	t := fmt.Sprintf("%d", time.Now().Unix())
	nonce := uuid.NewString()
	v := url.Values{
		"oauth_consumer_key":     {c.ConsumerKey},
		"oauth_nonce":            {nonce},
		"oauth_signature_method": {"HMAC-SHA1"},
		"oauth_timestamp":        {t},
	}
	baseString := fmt.Sprintf("%s&%s&%s", strings.ToUpper(r.Method),
		url.QueryEscape(r.URL.String()[:len(r.URL.String())-1]), url.QueryEscape(v.Encode()))
	mac := hmac.New(sha1.New, []byte(c.ConsumerSecret+"&"))
	mac.Write([]byte(baseString))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	headerVal := fmt.Sprintf(
		"OAuth oauth_signature=\"%s\","+
			"oauth_nonce=\"%s\","+
			"oauth_signature_method=\"HMAC-SHA1\","+
			"oauth_consumer_key=\"%s\","+
			"oauth_timestamp=\"%s\"", signature, nonce, c.ConsumerKey, t)
	r.Header.Set("Authorization", headerVal)

	return nil
}
