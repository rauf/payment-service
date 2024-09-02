package router

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
)

type Router struct {
	registry *gateway.Registry
}

func NewRouter(registry *gateway.Registry) *Router {
	return &Router{
		registry: registry,
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
		var result models.TransactionResponse
		slog.InfoContext(ctx, "Sending request to gateway", "gateway", g.Name())
		result, err = operation(g)
		if err == nil {
			return Response{
				Gateway: g.Name(),
				Data:    result,
			}, nil
		}
		if errors.Is(err, gateway.ErrGatewayUnavailable) {
			continue
		}
	}
	return Response{}, fmt.Errorf("all gateways failed: %w", err)
}
