package gateway

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rauf/payment-service/internal/backoff"
	"github.com/rauf/payment-service/internal/protocol"
	"github.com/rauf/payment-service/internal/serde"
)

// baseGateway is the base struct for all gateways.
type baseGateway[Req, Res any] struct {
	name            string
	serde           serde.Serde
	protocolHandler protocol.Handler
	retryConfig     backoff.RetryConfig
}

func newBaseGateway[Req, Res any](
	name string,
	serde serde.Serde,
	protocolHandler protocol.Handler,
	retryConfig backoff.RetryConfig,
) baseGateway[Req, Res] {
	return baseGateway[Req, Res]{
		name:            name,
		serde:           serde,
		protocolHandler: protocolHandler,
		retryConfig:     retryConfig,
	}
}

func (g *baseGateway[Req, Res]) sendWithRetry(ctx context.Context, data Req) (Res, error) {
	var zero Res
	var err error

	for attempt := 0; attempt <= g.retryConfig.MaxRetries; attempt++ {
		var response Res
		response, err = g.send(ctx, data)
		if err == nil {
			return response, nil
		}

		if errors.Is(err, ErrGatewayUnavailable) {
			slog.WarnContext(ctx, "Gateway unavailable, retrying", "attempt", attempt+1, "maxRetries", g.retryConfig.MaxRetries)
		} else if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return zero, fmt.Errorf("operation cancelled or timed out: %w", err)
		}
		if attempt == g.retryConfig.MaxRetries {
			return zero, fmt.Errorf("max retries reached, last error: %w", err)
		}

		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-time.After(g.retryConfig.Backoff.NextBackoff(attempt)):
		}
	}
	return zero, fmt.Errorf("all retries failed, last error: %w", err)
}

func (g *baseGateway[Req, Res]) send(ctx context.Context, data Req) (Res, error) {
	var zero Res
	if g.serde == nil {
		return zero, fmt.Errorf("data format is not initialized")
	}
	if g.protocolHandler == nil {
		return zero, fmt.Errorf("protocol handler is not initialized")
	}
	var buf bytes.Buffer
	err := g.serde.Serialize(&buf, data)
	if err != nil {
		return zero, fmt.Errorf("error marshaling data: %w", err)
	}

	response, err := g.protocolHandler.Send(ctx, buf.Bytes())
	if err != nil {
		return zero, fmt.Errorf("error sending data: %w", err)
	}
	if response == nil {
		return zero, fmt.Errorf("received nil response")
	}

	var result Res
	if err := g.serde.Deserialize(bytes.NewReader(response), &result); err != nil {
		return zero, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return result, nil
}

func (g *baseGateway[Req, Res]) Name() string {
	if g.name == "" {
		return "unnamed gateway"
	}
	return g.name
}
