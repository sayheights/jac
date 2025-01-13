package multipart

import (
	_ "embed"
	"testing"

	"github.com/google/uuid"
)

type testMarshal struct {
	Name       string          `multi:"name"`
	Countries  []string        `multi:"country"`
	UploadFile testMarshalJSON `multi:"uploadFile,application/json"`
}

type testMarshalJSON struct {
	ID      string    `json:"id"`
	Count   int       `json:"count"`
	Metrics []float64 `json:"metrics"`
}

var marshalJSONInput = testMarshalJSON{
	ID:      "ID_123",
	Count:   123,
	Metrics: []float64{1.3, 1.4, 1.6},
}

var testMarshalInput = testMarshal{
	Name:       "NAME_13",
	Countries:  []string{"us", "ca", "ad"},
	UploadFile: marshalJSONInput,
}

func TestMarshal(t *testing.T) {
	uuid.SetRand(nil)
	tests := []struct {
		name    string
		i       interface{}
		wantErr bool
	}{
		{
			i: testMarshalInput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Marshal(tt.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
