package multipart

import (
	"mime/multipart"
	"reflect"
	"testing"
)

func TestFieldKindError_Error(t *testing.T) {
	tests := []struct {
		name string
		kind string
		want string
	}{
		{
			kind: reflect.String.String(),
			want: "multipart: unhandled field kind string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FieldKindError{
				kind: tt.kind,
			}
			if got := f.Error(); got != tt.want {
				t.Errorf("FieldKindError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_basicFieldEncoder(t *testing.T) {
	type args struct {
		p  *formPart
		pw *multipart.Writer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := basicFieldEncoder(tt.args.p, tt.args.pw); (err != nil) != tt.wantErr {
				t.Errorf("basicFieldEncoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_arrayFieldEncoder(t *testing.T) {
	type args struct {
		p  *formPart
		pw *multipart.Writer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := arrayFieldEncoder(tt.args.p, tt.args.pw); (err != nil) != tt.wantErr {
				t.Errorf("arrayFieldEncoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_intEncoder(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intEncoder(tt.args.v); got != tt.want {
				t.Errorf("intEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_floatEncoder(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := floatEncoder(tt.args.v); got != tt.want {
				t.Errorf("floatEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boolEncoder(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := boolEncoder(tt.args.v); got != tt.want {
				t.Errorf("boolEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringEncoder(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringEncoder(tt.args.v); got != tt.want {
				t.Errorf("stringEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}
