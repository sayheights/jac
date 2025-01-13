package jac

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLinearBackoff(t *testing.T) {
	name := "LinearBackOff(%d)(%d) = got: %d, want: %d"
	tests := []struct {
		base time.Duration
		cnt  int
		want time.Duration
	}{
		{time.Second * 5, 2, time.Second * 10},
		{time.Minute * 7, 0, 0},
		{0, 5, time.Second * 5},
		{-1, 3, time.Second * 3},
		{-1, 0, 0},
	}

	for _, tt := range tests {
		got := LinearBackoff(tt.base)(tt.cnt)
		if got != tt.want {
			t.Fatalf(name, tt.base, tt.cnt, got, tt.want)
		}
	}
}

func TestLinearJitterBackoff(t *testing.T) {
	name := "LinearJitterBackOff(%d, %d, %d)(%d) = got: %v, want: %v"
	tests := []struct {
		min  time.Duration
		max  time.Duration
		seed int64
		cnt  int
		want time.Duration
	}{
		{cnt: 0, min: 1, max: 5, seed: 0, want: 0},
		{cnt: 1, min: 2, max: 5, seed: 0, want: 4},
		{cnt: 1, min: 5, max: 2, seed: 0, want: 4},
	}

	for _, tt := range tests {
		got := LinearJitterBackoff(tt.min, tt.max, tt.seed)(tt.cnt)
		if got != tt.want {
			t.Fatalf(name, tt.min, tt.max, tt.seed, tt.cnt, got, tt.want)
		}
	}
}

func TestLinearJitterBackoff_randomness(t *testing.T) {
	backoff := LinearJitterBackoff(time.Second*3, time.Second*120, 1000)
	var isRandom bool
	for i := 0; i < 100; i++ {
		first := backoff(1)
		second := backoff(1)
		if first != second {
			isRandom = true
			break
		}

	}

	assert.True(t, isRandom)
}

func TestExponentialBackoff(t *testing.T) {
	name := "ExponentialBackoff(%d, %d)(%d) = got: %d, want: %d"
	tests := []struct {
		min  time.Duration
		max  time.Duration
		cnt  int
		want time.Duration
	}{
		{time.Second, 5 * time.Minute, 0, time.Second},
		{time.Second, 5 * time.Minute, 1, 2 * time.Second},
		{time.Second, 5 * time.Minute, 2, 4 * time.Second},
		{time.Second, 5 * time.Minute, 3, 8 * time.Second},
		{time.Second, 5 * time.Minute, 63, 5 * time.Minute},
		{time.Second, 5 * time.Minute, 128, 5 * time.Minute},
	}

	for _, tt := range tests {
		got := ExponentialBackoff(tt.min, tt.max)(tt.cnt)
		if got != tt.want {
			t.Fatalf(name, tt.min, tt.max, tt.cnt, got, tt.want)
		}
	}
}
