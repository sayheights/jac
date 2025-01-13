package httptest

import (
	"reflect"
	"testing"
)

func Test_fromMap(t *testing.T) {
	type testCase struct {
		m    map[string][]string
		want []parameter
	}

	tests := []testCase{
		{
			m: map[string][]string{
				"z_key":     {"elem"},
				"key":       {"second", "first"},
				"first_key": {"3", "1", "2"},
			},
			want: []parameter{
				{key: "first_key", val: map[string]bool{"1": true, "2": true, "3": true}},
				{key: "key", val: map[string]bool{"first": true, "second": true}},
				{key: "z_key", val: map[string]bool{"elem": true}},
			},
		},
	}

	for i, tt := range tests {
		if got := fromMap(tt.m, ""); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d fromQuery() = %v, want %v", i, got, tt.want)
		}
	}
}

func Test_parameter_diff(t *testing.T) {
	type testCase struct {
		input []parameter
		want  *diff
	}

	tests := []testCase{
		{
			input: []parameter{{key: "no_key", val: map[string]bool{}}},
			want:  &diff{typ: missingKey, key: "key", val: "param,param1,param2"},
		},
		{
			input: []parameter{{key: "key", val: map[string]bool{"param1": true, "param2": true}}},
			want:  &diff{typ: missingValue, key: "key", val: "param"},
		},
		{
			input: []parameter{{key: "key", val: map[string]bool{
				"param":  true,
				"param1": true,
				"param2": true,
			}}},
			want: nil,
		},
		{
			input: []parameter{{key: "key", val: map[string]bool{
				"param":  true,
				"param1": true,
				"param2": true,
				"param3": true,
			}}},
			want: &diff{typ: unknownValue, key: "key", val: "param3"},
		},
	}

	for i, tt := range tests {
		p := parameter{
			key: "key",
			val: map[string]bool{
				"param":  true,
				"param1": true,
				"param2": true,
			},
			in: "TEST",
		}
		if tt.want != nil {
			tt.want.in = "TEST"
		}

		got := p.diff(tt.input)
		if got.report() != tt.want.report() {
			t.Errorf("%d: parameter.diff() = %s, want %s",
				i, got.report(), tt.want.report())
		}
	}
}

func Test_fromPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []parameter
	}{
		{
			in: "/pvar/pvar1/pvar2/",
			want: []parameter{
				{
					key: "SEGMENT_0",
					val: map[string]bool{"pvar": true},
					in:  "path",
				},
				{
					key: "SEGMENT_1",
					val: map[string]bool{"pvar1": true},
					in:  "path",
				},
				{
					key: "SEGMENT_2",
					val: map[string]bool{"pvar2": true},
					in:  "path",
				},
			},
		},
	}
	for _, tt := range tests {
		if got := fromPath(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("fromPath() = %v, want %v", got, tt.want)
		}
	}
}
