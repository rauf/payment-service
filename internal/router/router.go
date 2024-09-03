package router

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/registry"
	"github.com/sony/gobreaker/v2"
)

// Router is a struct that routes the request to the available gateways
// It has a registry of all available gateways and uses circuit breakers to prevent cascading failures.
type Router struct {
	*circuitBreakers
	registry *registry.Registry[gateway.PaymentGateway]
}

func NewRouter(
	registry *registry.Registry[gateway.PaymentGateway],
	settings gobreaker.Settings,
) *Router {
	return &Router{
		registry:        registry,
		circuitBreakers: newCircuitBreakers(settings),
	}
}

type Response struct {
	Gateway string
	Data    models.TransactionResponse
}

func (r *Router) SendMessage(ctx context.Context, preferredGateway string, operation func(gateway.PaymentGateway) (models.TransactionResponse, error)) (Response, error) {
	allGateways, err := r.registry.ListWithPreference(preferredGateway)
	if err != nil {
		return Response{}, fmt.Errorf("failed to get preferred gateways list: %w", err)
	}

	for _, g := range allGateways {
		done, cbErr := r.isRequestAllowed(ctx, g.Name())
		if cbErr != nil {
			continue
		}

		slog.InfoContext(ctx, "Sending request to gateway", "gateway", g.Name())

		var result models.TransactionResponse
		result, err = operation(g)
		if err == nil {
			done(true)
			return Response{
				Gateway: g.Name(),
				Data:    result,
			}, nil
		}

		done(false)
		if errors.Is(err, gateway.ErrGatewayUnavailable) {
			continue
		}
	}
	return Response{}, fmt.Errorf("all gateways failed: %w", err)
}
