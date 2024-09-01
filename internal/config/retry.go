package config

import "github.com/rauf/payment-service/internal/backoff"

type RetryConfig struct {
	MaxRetries int
	Backoff    backoff.Strategy
}
