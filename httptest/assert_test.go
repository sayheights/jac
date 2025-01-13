package httptest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAssertHeaders(t *testing.T) {
	type testCase struct {
		in   http.Header
		want int
	}
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	extraHeaders := http.Header{}
	extraHeaders.Add("Content-Type", "application/json")
	extraHeaders.Add("Accept", "gzip")

	tests := []testCase{
		{in: headers, want: 200},
		{in: http.Header{}, want: 400},
		{in: http.Header{"Content-Type": []string{"application/json", "text/csv"}}, want: 400},
		{in: extraHeaders, want: 200},
	}

	for i, tt := range tests {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header = tt.in

		rr := httptest.NewRecorder()
		handler := AssertHeaders(headers)(SuccessHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != tt.want {
			t.Errorf("%d: handler returned wrong status code: got %v want %v\n%s",
				i, rr.Code, tt.want, rr.Body.String())
		}
	}
}

func TestAssertQuery(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	query := "key=param&key1=param1&keys2=param2,param3"

	tests := []testCase{
		{in: query, want: 200},
		{in: "a=b", want: 400},
		{in: "key1=param1&keys2=param2,param3&key=param", want: 200},
		{in: query + "&extra=extra1", want: 400},
	}

	for i, tt := range tests {
		req, err := http.NewRequest("GET", "?"+tt.in, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := AssertQuery(query)(SuccessHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != tt.want {
			t.Errorf("%d: handler returned wrong status code: got %v want %v\n%s",
				i, rr.Code, tt.want, rr.Body.String())
		}
	}
}

func TestAssertPath(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	path := "/pathvar/pathvar1/pathvar2"

	tests := []testCase{
		{in: "pathvar1/pathvar/pathvar2/", want: 400},
		{in: path, want: 200},
		{in: "/pathvar/pavar2/pathvar3", want: 400},
		{in: "/pathvar/pathvar1/pathvar2/pathvar3/", want: 400},
	}

	for i, tt := range tests {
		req, err := http.NewRequest("GET", tt.in, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := AssertPath(path)(SuccessHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != tt.want {
			t.Errorf("%d: handler returned wrong status code: got %v want %v\n%s",
				i, rr.Code, tt.want, rr.Body.String())
		}

	}
}

func TestAssertMethod(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	method := "GET"

	tests := []testCase{
		{in: method, want: 200},
		{in: "POST", want: 400},
	}

	for i, tt := range tests {
		req, err := http.NewRequest(tt.in, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := AssertMethod(method)(SuccessHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != tt.want {
			t.Errorf("%d: handler returned wrong status code: got %v want %v\n%s",
				i, rr.Code, tt.want, rr.Body.String())
		}
	}
}

func TestAssertBody(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	body := `{"ints":[1,2,3],"string":"string_key","bool":true}`
	bodyNewLine := `{"":[1,2,3],"string":"string_key","bool":true, "cool":true}`

	tests := []testCase{
		{in: body, want: 200},
		{in: `{"ints":[1,2,3]}`, want: 400},
		{in: bodyNewLine, want: 400},
	}

	for i, tt := range tests {
		req, err := http.NewRequest("POST", "/", strings.NewReader(tt.in))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := AssertBody([]byte(body))(SuccessHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != tt.want {
			t.Errorf("%d: handler returned wrong status code: got %v want %v\n%s",
				i, rr.Code, tt.want, rr.Body.String())
		}
	}
}
