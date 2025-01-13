package httptest

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type Assert func(next http.Handler) http.Handler

func AssertMethod(method string) Assert {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				diff := &diff{
					typ: valueMismatch,
					key: method,
					in:  "method",
					val: r.Method,
				}
				diff.report()
				w.WriteHeader(400)
				fmt.Fprint(w, "Method Assertion Failed")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AssertPath(path string) Assert {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expected, actual := fromPath(path), fromPath(r.URL.Path)
			diff := diffParams(expected, actual, false)
			if diff != nil {
				diff.report()
				w.WriteHeader(400)
				fmt.Fprintln(w, "Path Assertion Failed")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AssertQuery(query string) Assert {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q, _ := url.ParseQuery(query)
			expected, actual := fromMap(q, "query"), fromMap(r.URL.Query(), "query")
			diff := diffParams(expected, actual, false)
			if diff != nil {
				diff.report()
				w.WriteHeader(400)
				fmt.Fprintf(w, "Query Assertion Failed")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AssertHeaders(headers http.Header) Assert {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expected, actual := fromMap(headers, "header"), fromMap(r.Header, "header")
			diff := diffParams(expected, actual, true)
			if diff != nil {
				diff.report()
				w.WriteHeader(400)
				fmt.Fprintln(w, "Header Assertion Failed")
				return
			}
			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func AssertBody(expected []byte) Assert {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actual, err := io.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(string(expected), string(actual), false)
			if !checkBodyEqual(actual, expected) {
				w.WriteHeader(400)
				fmt.Fprintln(w, "Body Assertion Failed")
				val := dmp.DiffPrettyText(diffs)
				diff := &diff{in: "body", typ: valueMismatch, key: "diff", val: val}
				diff.report()
				return
			}
			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func checkBodyEqual(got, want []byte) bool {
	sGot := strings.TrimSpace(string(got))
	sWant := strings.TrimSpace(string(want))
	return sGot == sWant
}
