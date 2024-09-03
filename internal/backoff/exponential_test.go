package backoff

import (
	"testing"
	"time"
)

func TestNewExponentialBackoff(t *testing.T) {
	initialInterval := time.Second
	multiplier := 2.0
	maxInterval := time.Minute

	eb := NewExponentialBackoff(initialInterval, multiplier, maxInterval)

	if eb.InitialInterval != initialInterval {
		t.Errorf("Expected InitialInterval to be %v, got %v", initialInterval, eb.InitialInterval)
	}
	if eb.Multiplier != multiplier {
		t.Errorf("Expected Multiplier to be %v, got %v", multiplier, eb.Multiplier)
	}
	if eb.MaxInterval != maxInterval {
		t.Errorf("Expected MaxInterval to be %v, got %v", maxInterval, eb.MaxInterval)
	}
}

func TestExponentialBackoff_NextBackoff(t *testing.T) {
	tests := []struct {
		name     string
		backoff  ExponentialBackoff
		attempts []int
		expected []time.Duration
	}{
		{
			name:     "Standard exponential backoff",
			backoff:  NewExponentialBackoff(time.Second, 2.0, time.Minute),
			attempts: []int{0, 1, 2, 3, 4},
			expected: []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second},
		},
		{
			name:     "Backoff with max interval",
			backoff:  NewExponentialBackoff(time.Second, 2.0, 10*time.Second),
			attempts: []int{0, 1, 2, 3, 4},
			expected: []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 10 * time.Second},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, attempt := range tt.attempts {
				got := tt.backoff.NextBackoff(attempt)
				if got != tt.expected[i] {
					t.Errorf("Attempt %d: expected %v, got %v", attempt, tt.expected[i], got)
				}
			}
		})
	}
}
