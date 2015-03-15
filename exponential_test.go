package backoff

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

const (
	// Multiplying expected by 1.4 yields consistently passing tests for me locally, anything less lead to transient failures
	// It's enough to prove the backoffs are at least almost exponential, while allowing for some timing variences between runs
	TIMING_JITTER_COMP float64 = float64(1.4)
)

func TestExponentialBackoff(t *testing.T) {
	exp := NewExponential(1*time.Millisecond, 0)
	for i := 0; i < 5; i++ {
		start := time.Now()
		exp.Backoff()
		end := time.Now()
		assert.InDelta(t, end.UnixNano(), start.UnixNano(), (TIMING_JITTER_COMP*math.Pow(2, float64(i)))*float64(time.Millisecond))
	}
}

func TestExponentialLimit(t *testing.T) {
	exp := NewExponential(250*time.Millisecond, time.Second)
	for i := 0; i < 6; i++ {
		exp.Backoff()
	}
	assert.Equal(t, 6, exp.count)
	assert.Equal(t, time.Second, exp.getBackoff())
}
