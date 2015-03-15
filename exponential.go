package backoff

import (
	"math"
	"time"
)

type Exponential struct {
	LastWait time.Duration
	start    time.Duration
	limit    time.Duration
	count    int
}

func (exp *Exponential) Reset() {
	exp.count = 0
}

func (exp *Exponential) getBackoff() time.Duration {
	period := int(math.Pow(2, float64(exp.count)))
	backoff := time.Duration(period) * exp.start
	if exp.limit > 0 && backoff > exp.limit {
		backoff = exp.limit
	}
	return backoff
}

func (exp *Exponential) Backoff() {
	backoff := exp.getBackoff()
	time.Sleep(backoff)
	exp.count++
	exp.LastWait = backoff
}

func NewExponential(start time.Duration, limit time.Duration) *Exponential {
	return &Exponential{start: start, limit: limit}
}
