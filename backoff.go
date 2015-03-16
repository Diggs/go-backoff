package backoff

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

// BackoffStrategy can be implemented to provide different backoff algorithms.
type BackoffStrategy interface {
	// GetBackoffDuration calculates the next time.Duration that the current thread will sleep for when backing off.
	// It receives the current backoff count, the initial backoff duration and the last back off duration.
	GetBackoffDuration(int, time.Duration, time.Duration) time.Duration
}

// Backoff tracks the generic state of the configured back off strategy.
type Backoff struct {
	// LastDuration contains the duration that was previously waited, or 0 if no backoff has occurred yet.
	LastDuration time.Duration
	// NextDuration contains the duration that will be waited on the next call to Backoff().
	NextDuration time.Duration
	start        time.Duration
	limit        time.Duration
	count        int
	strategy     BackoffStrategy
}

// Creates a new backoff using the specified BackoffStrategy, start duration and limit.
func NewBackoff(strategy BackoffStrategy, start time.Duration, limit time.Duration) *Backoff {
	backoff := Backoff{strategy: strategy, start: start, limit: limit}
	backoff.Reset()
	return &backoff
}

// Reset sets the Backoff to its initial conditions ready to start over.
func (exp *Backoff) Reset() {
	exp.count = 0
	exp.LastDuration = 0
	exp.NextDuration = exp.getNextDuration()
}

// Backoff causes the current thread/routine to sleep for NextDuration.
func (exp *Backoff) Backoff() {
	time.Sleep(exp.NextDuration)
	exp.count++
	exp.LastDuration = exp.NextDuration
	exp.NextDuration = exp.getNextDuration()
}

func (exp *Backoff) getNextDuration() time.Duration {
	backoff := exp.strategy.GetBackoffDuration(exp.count, exp.start, exp.LastDuration)
	if exp.limit > 0 && backoff > exp.limit {
		backoff = exp.limit
	}
	return backoff
}

type exponential struct{}

func (exponential) GetBackoffDuration(backoffCount int, start time.Duration, lastBackoff time.Duration) time.Duration {
	period := int64(math.Pow(2, float64(backoffCount)))
	return time.Duration(period) * start
}

// NewExponential creates a new backoff using the exponential backoff algorithm.
func NewExponential(start time.Duration, limit time.Duration) *Backoff {
	return NewBackoff(exponential{}, start, limit)
}

type exponentialFullJitter struct{}

func (exponentialFullJitter) GetBackoffDuration(backoffCount int, start time.Duration, lastBackoff time.Duration) time.Duration {
	backoff := exponential{}.GetBackoffDuration(backoffCount, start, lastBackoff)
	jitter, _ := rand.Int(rand.Reader, big.NewInt(int64(backoff)))
	return time.Duration(jitter.Int64())
}

// NewExponentialFullJitter creates a new backoff using the exponential with full jitter backoff algorithm.
func NewExponentialFullJitter(start time.Duration, limit time.Duration) *Backoff {
	return NewBackoff(exponentialFullJitter{}, start, limit)
}

type linear struct{}

func (linear) GetBackoffDuration(backoffCount int, start time.Duration, lastBackoff time.Duration) time.Duration {
	return time.Duration(backoffCount) * start
}

// NewLinear creates a new backoff using the linear backoff algorithm.
func NewLinear(start time.Duration, limit time.Duration) *Backoff {
	return NewBackoff(linear{}, start, limit)
}
