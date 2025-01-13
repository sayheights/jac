package jac

import "testing"

func TestMethod_String(t *testing.T) {
	tests := []struct {
		name string
		m    Method
		want string
	}{
		{
			m:    PATCH,
			want: "PATCH",
		},
		{
			m:    POST,
			want: "POST",
		},
		{
			m:    GET,
			want: "GET",
		},
		{
			m:    PUT,
			want: "PUT",
		},
		{
			m:    Method(1000),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Method.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMethod(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Method
	}{
		{
			in:   "Get",
			want: GET,
		},
		{
			in:   "POST",
			want: POST,
		},
		{
			in:   "UNKNOWN",
			want: GET,
		},
		{
			in:   "",
			want: GET,
		},
		{
			in:   "PuT",
			want: PUT,
		},
		{
			in:   "PAtch",
			want: PATCH,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMethod(tt.in); got != tt.want {
				t.Errorf("NewMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMethod_IsIdempotent(t *testing.T) {
	tests := []struct {
		name string
		m    Method
		want bool
	}{
		{
			m:    POST,
			want: false,
		},
		{
			m:    GET,
			want: true,
		},
		{
			m:    PUT,
			want: true,
		},
		{
			m:    DELETE,
			want: true,
		},
		{
			m:    HEAD,
			want: true,
		},
		{
			m:    PATCH,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsIdempotent(); got != tt.want {
				t.Errorf("Method.IsIdempotent() = %v, want %v", got, tt.want)
			}
		})
	}
}
