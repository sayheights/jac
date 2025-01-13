package jac

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_logger_Log(t *testing.T) {
	tests := []struct {
		name string
		t    *transaction
	}{
		{
			t: &transaction{
				state: txnExhausted,
				id:    "test-123",
			},
		},
	}
	for _, tt := range tests {
		testOut, _ := os.Create("test.log")
		os.Stderr = testOut
		t.Run(tt.name, func(t *testing.T) {
			l := newLogger("test.logger.com", "TestAPI")
			l.Log(tt.t)
			assert.FileExists(t, "test.log")
			os.Remove("test.log")
		})
	}
}
