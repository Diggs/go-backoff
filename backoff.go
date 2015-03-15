package backoff

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

type BackoffStrategy interface {
	GetBackoffDuration(int, time.Duration, time.Duration) time.Duration
}

type Backoff struct {
	LastBackoff time.Duration
	start       time.Duration
	limit       time.Duration
	count       int
	strategy    BackoffStrategy
}

func (exp *Backoff) Reset() {
	exp.count = 0
}

func (exp *Backoff) GetBackoffDuration() time.Duration {
	backoff := exp.strategy.GetBackoffDuration(exp.count, exp.start, exp.LastBackoff)
	if exp.limit > 0 && backoff > exp.limit {
		backoff = exp.limit
	}
	return backoff
}

func (exp *Backoff) Backoff() {
	backoff := exp.GetBackoffDuration()
	time.Sleep(backoff)
	exp.count++
	exp.LastBackoff = backoff
}

type exponential struct{}

func (exponential) GetBackoffDuration(backoffCount int, start time.Duration, lastBackoff time.Duration) time.Duration {
	period := int(math.Pow(2, float64(backoffCount)))
	return time.Duration(period) * start
}

func NewExponential(start time.Duration, limit time.Duration) *Backoff {
	return &Backoff{strategy: exponential{}, start: start, limit: limit}
}

type exponentialFullJitter struct{}

func (exponentialFullJitter) GetBackoffDuration(backoffCount int, start time.Duration, lastBackoff time.Duration) time.Duration {
	backoff := exponential{}.GetBackoffDuration(backoffCount, start, lastBackoff)
	jitter, _ := rand.Int(rand.Reader, big.NewInt(int64(backoff)))
	return time.Duration(jitter.Int64())
}

func NewExponentialFullJitter(start time.Duration, limit time.Duration) *Backoff {
	return &Backoff{strategy: exponentialFullJitter{}, start: start, limit: limit}
}

type linear struct{}

func (linear) GetBackoffDuration(backoffCount int, start time.Duration, lastBackoff time.Duration) time.Duration {
	return time.Duration(backoffCount) * start
}

func NewLinear(start time.Duration, limit time.Duration) *Backoff {
	return &Backoff{strategy: linear{}, start: start, limit: limit}
}
