package gateway

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rauf/payment-service/internal/config"
	"github.com/rauf/payment-service/internal/format"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/protocol"
)

var (
	ErrGatewayUnavailable = fmt.Errorf("gateway unavailable")
	ErrorGatewayNotFound  = fmt.Errorf("gateway not found")
)

type PaymentGateway interface {
	Name() string
	Deposit(context.Context, models.DepositRequest) (models.DepositResponse, error)
	Withdraw(context.Context, models.WithdrawalRequest) (models.WithdrawalResponse, error)
}

type baseGateway struct {
	name            string
	dataFormat      format.DataFormat
	protocolHandler protocol.Handler
	retryConfig     config.RetryConfig
}

func newBaseGateway(name string, dataFormat format.DataFormat, protocolHandler protocol.Handler, retryConfig config.RetryConfig) baseGateway {
	return baseGateway{
		name:            name,
		dataFormat:      dataFormat,
		protocolHandler: protocolHandler,
		retryConfig:     retryConfig,
	}
}

func (g *baseGateway) SendWithRetry(ctx context.Context, data any) (any, error) {
	var lastError error
	for attempt := 0; attempt <= g.retryConfig.MaxRetries; attempt++ {
		response, err := g.Send(ctx, data)
		if err == nil {
			return response, nil
		}
		lastError = err

		if errors.Is(err, ErrGatewayUnavailable) {
			slog.WarnContext(ctx, "Gateway unavailable, retrying", "attempt", attempt+1, "maxRetries", g.retryConfig.MaxRetries)
		} else if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("operation cancelled or timed out: %w", err)
		}
		if attempt == g.retryConfig.MaxRetries {
			return nil, fmt.Errorf("max retries reached, last error: %w", lastError)
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-time.After(g.retryConfig.Backoff.NextBackoff(attempt)):
		}
	}
	return nil, fmt.Errorf("all retries failed, last error: %w", lastError)
}

func (g *baseGateway) Send(ctx context.Context, data any) (any, error) {
	if g.dataFormat == nil {
		return nil, fmt.Errorf("data format is not initialized")
	}
	if g.protocolHandler == nil {
		return nil, fmt.Errorf("protocol handler is not initialized")
	}
	encoded, err := g.dataFormat.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %w", err)
	}

	response, err := g.protocolHandler.Send(ctx, encoded)
	if err != nil {
		return nil, fmt.Errorf("error sending data: %w", err)
	}
	if response == nil {
		return nil, fmt.Errorf("received nil response")
	}

	var result any
	if err := g.dataFormat.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return result, nil
}

func (g *baseGateway) Name() string {
	if g.name == "" {
		return "unnamed gateway"
	}
	return g.name
}
