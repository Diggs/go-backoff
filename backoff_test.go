package backoff

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	exp := NewExponential(time.Millisecond, 0)
	for i := 0; i < 5; i++ {
		assert.Equal(t, int64(float64(time.Millisecond) * math.Pow(2, float64(i))), exp.GetBackoffDuration())
		exp.Backoff()
	}
}

func TestExponentialLimit(t *testing.T) {
	exp := NewExponential(250*time.Millisecond, time.Second)
	for i := 0; i < 5; i++ {
		exp.Backoff()
	}
	assert.Equal(t, 5, exp.count)
	assert.Equal(t, time.Second, exp.GetBackoffDuration())
}

func TestLinearBackoff(t *testing.T) {
	lin := NewLinear(time.Millisecond, 0)
	for i := 0; i < 5; i++ {
		assert.Equal(t, int64(float64(time.Millisecond) * float64(i)), lin.GetBackoffDuration())
		lin.Backoff()
	}
}

func TestLinearLimit(t *testing.T) {
	lin := NewLinear(250*time.Millisecond, time.Second)
	for i := 0; i < 5; i++ {
		lin.Backoff()
	}
	assert.Equal(t, 5, lin.count)
	assert.Equal(t, time.Second, lin.GetBackoffDuration())
}
