package multipart

import "testing"

func TestContentType_String(t *testing.T) {
	tests := []struct {
		name string
		c    ContentType
		want string
	}{
		{
			c:    TextCSS,
			want: "text/css",
		},
		{
			c:    ApplicationJSON,
			want: "application/json",
		},
		{
			c:    TextCSV,
			want: "text/csv",
		},
		{
			c:    ContentType(10239),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("ContentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
