package multipart

import (
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    tag
		wantErr bool
	}{
		{
			name: "tag",
			in:   `multi:"name,application/json"`,
			want: tag{fieldName: "name", contentType: ApplicationJSON},
		},
		{
			name: "tag",
			in:   `multi:"name,  text/csv  "`,
			want: tag{fieldName: "name", contentType: TextCSV},
		},
		{
			name: "tag",
			in:   `multi:"name"`,
			want: tag{fieldName: "name"},
		},
		{
			name: "tag",
			in:   `multi:",text/csv"`,
			want: tag{contentType: TextCSV},
		},
		{
			name:    "tag",
			in:      `multi:",textcsv"`,
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTag(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("%d: parse() error = %v, wantErr %v",
					i, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%d: parse() = %v, want %v",
					i, got, tt.want)
			}
		})
	}
}
