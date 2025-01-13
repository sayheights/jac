package jac

import "strings"

// Method represents an HTTP Method.
type Method int

// String returns the string representation of an HTTP method.
func (m Method) String() string {
	return methodToStr[m]
}

const (
	// GET method.
	GET Method = 1 << iota

	// POST method.
	POST

	// PUT method.
	PUT

	// DELETE method.
	DELETE

	// PATCH method.
	PATCH

	// HEAD method.
	HEAD

	// OPTIONS method.
	OPTIONS

	// TRACE method.
	TRACE
)

// idempotentMethods is bitmask for idempotent HTTP methods.
const idempotentMethods = GET | PUT | DELETE | HEAD

// strToMethod is the mapping from a string to a Method.
var strToMethod = map[string]Method{
	"GET":     GET,
	"POST":    POST,
	"TRACE":   TRACE,
	"OPTIONS": OPTIONS,
	"HEAD":    HEAD,
	"PATCH":   PATCH,
	"DELETE":  DELETE,
	"PUT":     PUT,
}

// methodToStr is the mapping from a Method to its string representation.
var methodToStr = map[Method]string{
	GET:     "GET",
	POST:    "POST",
	TRACE:   "TRACE",
	OPTIONS: "OPTIONS",
	HEAD:    "HEAD",
	PATCH:   "PATCH",
	DELETE:  "DELETE",
	PUT:     "PUT",
}

// NewMethod returns an Method instance from the given string.
func NewMethod(s string) Method {
	s = strings.ToUpper(s)
	m, ok := strToMethod[s]
	if !ok {
		return GET
	}

	return m
}

// IsIdempotent returns true if the method is idempotent.
func (m Method) IsIdempotent() bool {
	return m&idempotentMethods != 0
}
