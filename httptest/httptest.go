package httptest

import (
	"net/http"
	"net/http/httptest"

	"github.com/darrae/jac"
)

func ServerAndClient(h http.Handler, asserts ...Assert) (*httptest.Server, *jac.Client) {
	next := h
	for _, assert := range asserts {
		next = assert(next)
	}
	svr := httptest.NewServer(next)
	c := &jac.Client{BaseURL: svr.URL, DisableLogging: true}
	return svr, c
}

var SuccessHandler = ResponseHandler([]byte("SUCCESS"))

func ResponseHandler(response []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(response)
	}
}
