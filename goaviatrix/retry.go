package goaviatrix

import "time"

// RetryConfig holds configuration for retry behavior.
type RetryConfig struct {
	// MaxTries is the total number of attempts (including the first).
	MaxTries int
	// Backoff is the initial sleep duration between attempts.
	Backoff time.Duration
	// BackoffFunc advances the backoff duration after each failed attempt.
	// Defaults to exponential doubling when nil.
	BackoffFunc func(time.Duration) time.Duration
	// ShouldRetry determines whether the error warrants another attempt.
	// Defaults to always retrying when nil.
	ShouldRetry func(error) bool
	// OnRetry is called before sleeping between attempts. Useful for logging.
	// Receives the current attempt number (1-based), the sleep duration, and the error.
	OnRetry func(attempt int, backoff time.Duration, err error)
}

// Retry executes fn up to cfg.MaxTries times, sleeping between failures
// according to cfg.Backoff and cfg.BackoffFunc. It returns immediately when
// fn succeeds or when ShouldRetry returns false for the error.
func Retry(cfg RetryConfig, fn func() error) error {
	backoffFn := cfg.BackoffFunc
	if backoffFn == nil {
		backoffFn = func(d time.Duration) time.Duration { return d * 2 }
	}

	backoff := cfg.Backoff
	var err error
	for try := 1; try <= cfg.MaxTries; try++ {
		err = fn()
		if err == nil {
			return nil
		}
		if cfg.ShouldRetry != nil && !cfg.ShouldRetry(err) {
			return err
		}
		if try == cfg.MaxTries {
			break
		}
		if cfg.OnRetry != nil {
			cfg.OnRetry(try, backoff, err)
		}
		time.Sleep(backoff)
		backoff = backoffFn(backoff)
	}
	return err
}
