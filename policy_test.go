package jac

import (
	"net/http"
	"testing"
)

func TestRetryOn(t *testing.T) {
	name := "RetryOn() = want: %s, got: %s"
	tests := []struct {
		codes []int
		res   *http.Response
		want  bool
	}{
		{
			codes: defRetryCodes,
			res:   &http.Response{StatusCode: 100},
			want:  false,
		},
		{
			codes: defRetryCodes,
			res:   &http.Response{StatusCode: 429},
			want:  true,
		},
		{
			codes: nil,
			res:   &http.Response{StatusCode: 200},
			want:  false,
		},
		{
			codes: []int{503},
			res:   &http.Response{StatusCode: 500},
			want:  false,
		},
	}
	for _, tt := range tests {
		got := RetryOn(tt.codes)(tt.res)
		if got != tt.want {
			t.Fatalf(name, tt.want, got)
		}
	}
}

func TestRetryIdempotentsOn(t *testing.T) {
	idempoReq := &http.Request{Method: "GET"}
	nonIdempoReq := &http.Request{Method: "POST"}
	name := "RetryOn() = want: %s, got: %s"
	tests := []struct {
		codes []int
		res   *http.Response
		want  bool
	}{
		{
			codes: defRetryCodes,
			res:   &http.Response{StatusCode: 100, Request: idempoReq},
			want:  false,
		},
		{
			codes: defRetryCodes,
			res:   &http.Response{StatusCode: 429, Request: nonIdempoReq},
			want:  false,
		},
		{
			codes: nil,
			res:   &http.Response{StatusCode: 429, Request: nonIdempoReq},
			want:  false,
		},
		{
			codes: []int{503},
			res:   &http.Response{StatusCode: 503, Request: idempoReq},
			want:  true,
		},
	}
	for _, tt := range tests {
		got := RetryIdempotentsOn(tt.codes)(tt.res)
		if got != tt.want {
			t.Fatalf(name, tt.want, got)
		}
	}
}
