package router

import (
	"context"
	"testing"
	"time"

	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCircuitBreakers(t *testing.T) {
	settings := gobreaker.Settings{
		Name: "test",
	}
	cbs := newCircuitBreakers(settings)

	assert.NotNil(t, cbs.circuitBreakers)
	assert.Equal(t, settings, cbs.settings)
}

func TestCircuitBreakers_IsRequestAllowed(t *testing.T) {
	settings := gobreaker.Settings{
		Name:     "test",
		Interval: 5 * time.Second,
		Timeout:  1 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	}
	cbs := newCircuitBreakers(settings)

	t.Run("request allowed", func(t *testing.T) {
		done, err := cbs.isRequestAllowed(context.Background(), "gateway1")
		require.NoError(t, err)
		assert.NotNil(t, done)

		done(true)
	})

	t.Run("circuit open", func(t *testing.T) {
		// Trigger the circuit breaker to open
		for i := 0; i < 4; i++ {
			done, _ := cbs.isRequestAllowed(context.Background(), "gateway1")
			done(false)
		}

		_, err := cbs.isRequestAllowed(context.Background(), "gateway1")
		assert.Error(t, err)
	})
}

func TestCircuitBreakers_GetCircuitBreaker(t *testing.T) {
	settings := gobreaker.Settings{
		Name: "test",
	}
	cbs := newCircuitBreakers(settings)

	t.Run("get existing circuit breaker", func(t *testing.T) {
		err := cbs.ensureCircuitBreaker("gateway1")
		require.NoError(t, err)

		cb, err := cbs.getCircuitBreaker("gateway1")
		assert.NoError(t, err)
		assert.NotNil(t, cb)
	})

	t.Run("get non-existing circuit breaker", func(t *testing.T) {
		cb, err := cbs.getCircuitBreaker("gateway2")
		assert.NoError(t, err)
		assert.NotNil(t, cb)
	})
}

func TestCircuitBreakers_EnsureCircuitBreaker(t *testing.T) {
	settings := gobreaker.Settings{
		Name: "test",
	}
	cbs := newCircuitBreakers(settings)

	err := cbs.ensureCircuitBreaker("gateway1")
	assert.NoError(t, err)

	cb, err := cbs.circuitBreakers.Get("gateway1")
	assert.NoError(t, err)
	assert.NotNil(t, cb)
}
