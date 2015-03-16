## What is it?

go-backoff implements linear, [exponential and exponential with full jitter](http://www.awsarchitectureblog.com/2015/03/backoff.html) back off algorithms.

## Usage

Instantiate and use a Backoff instance for the lifetime of a go routine. Call backoff.Backoff() to put the routine to sleep for the desired amount of time. Use backoff.Reset() to reset counters back to 0.

```go
  // Back off linearly, starting at 250ms, capping at 16 seconds
  linear = backoff.NewLinear(250*time.Millisecond, 16*time.Second)
  // Back off exponentially, starting at 5 seconds, capping at 320 seconds
  exp = backoff.NewExponential(5*time.Second, 320*time.Second)
  // Back off exponentially, starting at 1 minute, with no cap
  expt = backoff.NewExponential(time.Minute, 0)
  // Back off between 0 and exponentially, starting at 30 seconds, capping at 10 minutes
  expj = backoff.NewExponentialFullJitter(30*time.Second, 10*time.Minute)

  ...

  for {
    err := tryDoThing()
    if err != nil {
      glog.Debugf("Backing off %d second(s)", exp.NextDuration/time.Second)
      exp.Backoff()
    } else {
      exp.Reset()
      return
    }
  }

```

## Tests
```bash
go get github.com/tools/godep
godep restore
go test

PASS
ok    github.com/diggs/go-backoff 6.319s
```