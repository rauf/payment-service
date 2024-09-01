package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/rauf/payment-service/internal/config"
	"github.com/rauf/payment-service/internal/format"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/protocol"
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

func (g *baseGateway) SendWithRetry(ctx context.Context, data any) (any, error) {
	var lastError error
	for attempt := 0; attempt <= g.retryConfig.MaxRetries; attempt++ {
		response, err := g.Send(ctx, data)
		if err == nil {
			return response, nil
		}
		lastError = err

		if attempt == g.retryConfig.MaxRetries {
			return nil, fmt.Errorf("max retries reached, last error: %w", lastError)
		}

		// exponential backoff
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(g.retryConfig.Backoff.NextBackoff(attempt)):
		}
	}
	return nil, lastError
}

func (g *baseGateway) Send(ctx context.Context, data any) (any, error) {
	encoded, err := g.dataFormat.Marshal(data)
	if err != nil {
		return nil, err
	}

	response, err := g.protocolHandler.Send(ctx, encoded)

	var result any
	if err := g.dataFormat.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return result, nil
}

func (g *baseGateway) Name() string {
	return g.name
}
