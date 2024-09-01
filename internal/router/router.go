package router

import (
	"errors"
	"fmt"

	"github.com/rauf/payment-service/internal/gateway"
)

type Router struct {
	registry *gateway.Registry
}

func NewRouter(registry *gateway.Registry) *Router {
	return &Router{
		registry: registry,
	}
}

func (r *Router) SendMessage(preferredGateway string, operation func(gateway.PaymentGateway) (any, error)) (any, error) {
	allGateways, err := r.registry.ListWithPreference(preferredGateway)
	if err != nil {
		return "", fmt.Errorf("failed to get preferred gateways list: %w", err)
	}

	var lastError error
	for _, g := range allGateways {
		result, err := operation(g)
		if err == nil {
			return result, nil
		}
		lastError = err
		if errors.Is(err, gateway.ErrGatewayUnavailable) {
			continue
		}
	}
	return "", fmt.Errorf("all gateways failed: %w", lastError)
}
