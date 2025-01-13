package jac

import (
	"math"
	"math/rand"
	"time"
)

const (
	// defMinWait is the default minimum wait time
	// parameter used by the ExponentialBackoff algorithm.
	defMinWait = 1 * time.Second

	// defMaxWait is the default maximum wait time
	// parameter used by the ExponentialBackoff algorithm.
	defMaxWait = 30 * time.Second

	// defBaseWait is the default base wait time.
	defBaseWait = 1 * time.Second
)

// DefaultBackoff is the defaul implementation of ExponentialBackoff.
var DefaultBackoff Backoff = ExponentialBackoff(defMinWait, defMaxWait)

// Backoff determines the duration to wait before retrying a request.
//
// The duration is adjusted based on the amount of retries already attempted.
// Generally the wait duration increases in proportion with the retry
// count; however, no formal constraints are imposed on the implementations.
//
// Backoff only accepts the attempt count as the input and parameterization
// should be handled via function closures.
//
// The actual wait time can be shorter than the duration provided by
// the Backoff as an external cancellation or Deadline timeout of
// a Request's context would immediately return.
//
// In cases where the server response includes a Retry-After, Backoff is
// ignored and the provided duration is used.
type Backoff func(cnt int) time.Duration

// LinearBackoff is a Backoff implementation where the wait duration
// increases in a linear fashion as the current count is multiplied
// by the base duration. If the input base duration is smaller than
// or equal to 0, the default base wait duration of 1 second is used.
func LinearBackoff(base time.Duration) Backoff {
	if base <= 0 {
		base = defBaseWait
	}
	return func(cnt int) time.Duration {
		return base * time.Duration(cnt)
	}
}

// LinearJitterBackoff randomized its wait duration within the provided
// range as a precaution against a potential thundering herd problem.
// If the provided min value is larger than the maximum their values are swapped.
func LinearJitterBackoff(min, max time.Duration, seed int64) Backoff {
	rnd := rand.New(rand.NewSource(seed))
	return func(cnt int) time.Duration {
		if min > max {
			tmp := max
			max, min = min, tmp
		}
		jitter := rnd.Float64() * float64(max-min)
		jitterMin := int64(jitter) + int64(min)

		return time.Duration(jitterMin * int64(cnt))
	}
}

// ExponentialBackoff spaces out the wait duration exponentially as the
// attempt count increases while staying within the provided range.
func ExponentialBackoff(min, max time.Duration) Backoff {
	return func(cnt int) time.Duration {
		mult := math.Pow(2, float64(cnt)) * float64(min)
		wait := time.Duration(mult)
		if float64(wait) != mult || wait > max {
			wait = max
		}

		return wait
	}
}
