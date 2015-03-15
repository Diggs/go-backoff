package backoff

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLinearBackoff(t *testing.T) {
	exp := NewLinear(1*time.Millisecond, 0)
	for i := 0; i < 5; i++ {
		start := time.Now()
		exp.Backoff()
		end := time.Now()
		assert.InDelta(t, end.UnixNano(), start.UnixNano(), (TIMING_JITTER_COMP*float64(i))*float64(time.Millisecond))
	}
}

func TestLinearLimit(t *testing.T) {
	exp := NewLinear(250*time.Millisecond, time.Second)
	for i := 0; i < 6; i++ {
		exp.Backoff()
	}
	assert.Equal(t, 6, exp.count)
	assert.Equal(t, time.Second, exp.getBackoff())
}
