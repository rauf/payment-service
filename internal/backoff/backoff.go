package backoff

import (
	"time"
)

type Strategy interface {
	NextBackoff(attempt int) time.Duration
}

type RetryConfig struct {
	MaxRetries int
	Backoff    Strategy
}
