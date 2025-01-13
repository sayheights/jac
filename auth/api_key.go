package auth

import (
	"errors"
	"net/http"
)

type APIKey struct {
	Key   string
	Value string
	In    In
}

func NewAPIKey(key, value string, in In) *APIKey {
	return &APIKey{Key: key, Value: value, In: in}
}

// Authorize sets the api key value to url query or header of request with looking In value
func (a *APIKey) Authorize(r *http.Request) error {
	switch a.In {
	case InQuery:
		q := r.URL.Query()
		q.Set(a.Key, a.Value)
		r.URL.RawQuery = q.Encode()
		return nil
	case InHeader:
		r.Header.Set(a.Key, a.Value)
		return nil
	}

	return errors.New("auth: api key can not be in request body")
}
