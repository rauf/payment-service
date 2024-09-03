package backoff

import (
	"time"
)

// Strategy is an interface that needs to be implemented by any backoff strategy.
type Strategy interface {
	NextBackoff(attempt int) time.Duration
}

type RetryConfig struct {
	MaxRetries int
	Backoff    Strategy
}
