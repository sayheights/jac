package auth

import "net/http"

type Basic struct {
	ID  string
	Key string
}

func NewBasic(username, password string) *Basic {
	return &Basic{ID: username, Key: password}
}

// Authorize simply sets auth id and key to the request
func (b *Basic) Authorize(r *http.Request) error {
	r.SetBasicAuth(b.ID, b.Key)
	return nil
}
