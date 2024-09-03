package router

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/registry"
	"github.com/sony/gobreaker/v2"
)

// circuitBreakers is a registry of circuit breakers for all the payment gateways
type circuitBreakers struct {
	circuitBreakers *registry.Registry[*gobreaker.TwoStepCircuitBreaker[gateway.PaymentGateway]]
	settings        gobreaker.Settings
}

func newCircuitBreakers(settings gobreaker.Settings) *circuitBreakers {
	return &circuitBreakers{
		circuitBreakers: registry.NewRegistry[*gobreaker.TwoStepCircuitBreaker[gateway.PaymentGateway]](),
		settings:        settings,
	}
}

func (cbs *circuitBreakers) isRequestAllowed(ctx context.Context, gatewayName string) (func(success bool), error) {
	cb, cbErr := cbs.getCircuitBreaker(gatewayName)
	if cbErr != nil {
		return nil, fmt.Errorf("failed to get circuit breaker: %w", cbErr)
	}
	done, cbErr := cb.Allow()
	if cbErr != nil {
		slog.ErrorContext(ctx, "circuit is open for gateway", "error", cbErr, "gateway", gatewayName)
	}
	return done, cbErr
}

func (cbs *circuitBreakers) getCircuitBreaker(gatewayName string) (*gobreaker.TwoStepCircuitBreaker[gateway.PaymentGateway], error) {
	if err := cbs.ensureCircuitBreaker(gatewayName); err != nil {
		return nil, fmt.Errorf("failed to ensure circuit breaker: %w", err)
	}
	cb, _ := cbs.circuitBreakers.Get(gatewayName)
	return cb, nil
}

func (cbs *circuitBreakers) ensureCircuitBreaker(gatewayName string) error {
	if _, err := cbs.circuitBreakers.Get(gatewayName); err == nil {
		return nil
	}

	settings := cbs.settings
	settings.Name = gatewayName
	cb := gobreaker.NewTwoStepCircuitBreaker[gateway.PaymentGateway](settings)
	if err := cbs.circuitBreakers.Register(gatewayName, cb); err != nil {
		return fmt.Errorf("failed to register circuit breaker: %w", err)
	}
	return nil
}
