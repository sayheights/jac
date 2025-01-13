package jac

import (
	"net/http"
	"strings"
)

// defRetryCodes is the default retry code values
// used by the RetryOn policy.
var defRetryCodes = []int{429, 500, 502, 503, 504}

// DefaultPolicy is the default implementation of the RetryOn policy.
var DefaultPolicy RetryPolicy = RetryOn(defRetryCodes)

// RetryPolicy specifies the HTTP Responses that can be
// retried when encountered.It is called after every unsuccessful
// transaction and return true if it deems the situation recoverable.
// A request may not be a retried even when the RetryPolicy return
// true, for instance, when the maximum attempt count is reached.
type RetryPolicy func(*http.Response) bool

// RetryOn is a RetryPolicy where the request is only retried if the status
// code of the response is included in the input.
func RetryOn(codes []int) RetryPolicy {
	return func(res *http.Response) bool {
		for _, code := range codes {
			if code == res.StatusCode {
				return true
			}
		}
		return false
	}
}

// idempontentMethods specifies the HTTP methods that should behave
// idempotently if implemented correctly.
var idempontentMethods = map[string]bool{
	"GET":    true,
	"HEAD":   true,
	"PUT":    true,
	"DELETE": true,
}

// RetryIdempotentsOn policy returns false for all non-idempotent
// operations and true for idempotent requests whose status code
// is among the provided codes in the input.
//
// It behaves exactly as RetryOn policy for idempotent requests.
func RetryIdempotentsOn(codes []int) RetryPolicy {
	retryOn := RetryOn(codes)
	return func(res *http.Response) bool {
		method := strings.ToUpper(res.Request.Method)
		isIdempotent := idempontentMethods[method]
		if isIdempotent {
			return retryOn(res)
		}

		return false
	}
}
