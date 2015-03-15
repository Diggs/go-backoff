package backoff

import (
  "time"
)

type Linear struct {
  LastWait time.Duration
  increment time.Duration
  limit time.Duration
  count int
}

func(lin *Linear) Reset() {
  lin.count = 0
}

func(lin *Linear) getBackoff() time.Duration {
  backoff := time.Duration(lin.count) * lin.increment
  if lin.limit > 0 && backoff > lin.limit {
    backoff = lin.limit
  }
  return backoff
}

func(lin *Linear) Backoff() {
  backoff := lin.getBackoff()
  time.Sleep(backoff)
  lin.count++
  lin.LastWait = backoff
}

func NewLinear(increment time.Duration, limit time.Duration) *Linear {
  return &Linear{increment:increment, limit:limit}
}