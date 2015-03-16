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
		assert.Equal(t, int64(float64(time.Millisecond)*math.Pow(2, float64(i))), exp.NextDuration)
		exp.Backoff()
	}
}

func TestExponentialLimit(t *testing.T) {
	exp := NewExponential(time.Millisecond, 4*time.Millisecond)
	for i := 0; i < 5; i++ {
		exp.Backoff()
	}
	assert.Equal(t, 5, exp.count)
	assert.Equal(t, 4*time.Millisecond, exp.NextDuration)
}

func TestLinearBackoff(t *testing.T) {
	lin := NewLinear(time.Millisecond, 0)
	for i := 0; i < 5; i++ {
		assert.Equal(t, int64(float64(time.Millisecond)*float64(i)), lin.NextDuration)
		lin.Backoff()
	}
}

func TestLinearLimit(t *testing.T) {
	lin := NewLinear(time.Millisecond, 4*time.Millisecond)
	for i := 0; i < 5; i++ {
		lin.Backoff()
	}
	assert.Equal(t, 5, lin.count)
	assert.Equal(t, 4*time.Millisecond, lin.NextDuration)
}
