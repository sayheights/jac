package multipart

import (
	"reflect"
	"testing"
)

func Test_isNumericKind(t *testing.T) {
	tests := []struct {
		k    reflect.Kind
		want bool
	}{
		{
			k:    reflect.Int32,
			want: true,
		},
		{
			k:    reflect.Int8,
			want: true,
		},
		{
			k:    reflect.Float32,
			want: true,
		},
		{
			k:    reflect.Float64,
			want: true,
		},
		{
			k:    reflect.Array,
			want: false,
		},
	}
	for _, tt := range tests {
		if got := isNumericKind(tt.k); got != tt.want {
			t.Errorf("isNumericKind() = %v, want %v", got, tt.want)
		}
	}
}

func Test_isBasicKind(t *testing.T) {
	type args struct{}
	tests := []struct {
		k    reflect.Kind
		want bool
	}{
		{
			k:    reflect.String,
			want: true,
		},
		{
			k:    reflect.Int,
			want: true,
		},
		{
			k:    reflect.Struct,
			want: false,
		},
	}
	for _, tt := range tests {
		if got := isBasicKind(tt.k); got != tt.want {
			t.Errorf("isBasicKind() = %v, want %v", got, tt.want)
		}
	}
}

func Test_isIntKind(t *testing.T) {
	tests := []struct {
		k    reflect.Kind
		want bool
	}{
		{
			k:    reflect.Float32,
			want: false,
		},
		{
			k:    reflect.Float64,
			want: false,
		},
		{
			k:    reflect.Int64,
			want: true,
		},
	}
	for _, tt := range tests {
		if got := isIntKind(tt.k); got != tt.want {
			t.Errorf("isIntKind() = %v, want %v", got, tt.want)
		}
	}
}

func Test_isFloatKind(t *testing.T) {
	tests := []struct {
		k    reflect.Kind
		want bool
	}{
		{
			k:    reflect.Float32,
			want: true,
		},
		{
			k:    reflect.Float64,
			want: true,
		},
		{
			k:    reflect.Int64,
			want: false,
		},
	}

	for _, tt := range tests {
		if got := isFloatKind(tt.k); got != tt.want {
			t.Errorf("isFloatKind() = %v, want %v", got, tt.want)
		}
	}
}
