package backoff

import (
	"math"
	"time"
)

type Strategy interface {
	NextBackoff(attempt int) time.Duration
}

type ExponentialBackoff struct {
	InitialInterval time.Duration
	Multiplier      float64
	MaxInterval     time.Duration
}

func NewExponentialBackoff(initialInterval time.Duration, multiplier float64, maxInterval time.Duration) ExponentialBackoff {
	return ExponentialBackoff{
		InitialInterval: initialInterval,
		Multiplier:      multiplier,
		MaxInterval:     maxInterval,
	}
}

func (b ExponentialBackoff) NextBackoff(attempt int) time.Duration {
	interval := float64(b.InitialInterval) * math.Pow(b.Multiplier, float64(attempt))
	if interval > float64(b.MaxInterval) {
		interval = float64(b.MaxInterval)
	}
	return time.Duration(interval)
}
