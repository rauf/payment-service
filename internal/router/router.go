package router

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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

func (r *Router) SendMessage(ctx context.Context, preferredGateway string, operation func(gateway.PaymentGateway) (any, error)) (any, error) {
	allGateways, err := r.registry.ListWithPreference(preferredGateway)
	if err != nil {
		return "", fmt.Errorf("failed to get preferred gateways list: %w", err)
	}

	for _, g := range allGateways {
		var result any
		slog.InfoContext(ctx, "Sending request to gateway", "gateway", g.Name())
		result, err = operation(g)
		if err == nil {
			return result, nil
		}
		if errors.Is(err, gateway.ErrGatewayUnavailable) {
			continue
		}
	}
	return "", fmt.Errorf("all gateways failed: %w", err)
}
