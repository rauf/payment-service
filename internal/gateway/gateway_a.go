package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rauf/payment-service/internal/backoff"
	"github.com/rauf/payment-service/internal/config"
	"github.com/rauf/payment-service/internal/format"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/protocol"
)

type GatewayA struct {
	baseGateway
}

func NewGatewayA(address string) *GatewayA {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &GatewayA{
		baseGateway: baseGateway{
			name:            "gateway-a",
			dataFormat:      format.NewJSONProtocol(),
			protocolHandler: protocol.NewHTTPConnection(httpClient, http.MethodPost, address),
			retryConfig: config.RetryConfig{
				MaxRetries: 3,
				Backoff:    backoff.NewExponentialBackoff(1*time.Second, 1.2, 2*time.Second),
			},
		},
	}
}

func (g *GatewayA) Deposit(ctx context.Context, deposit models.DepositRequest) (models.DepositResponse, error) {
	request := map[string]any{
		"type":     "deposit",
		"amount":   deposit.Amount,
		"currency": deposit.Currency,
	}
	_, err := g.SendWithRetry(ctx, request)
	if err != nil {
		return models.DepositResponse{}, fmt.Errorf("error sending deposit request: %w", err)
	}

	return models.DepositResponse{
		TransactionID: "test",
	}, nil
}

func (g *GatewayA) Withdraw(ctx context.Context, withdrawal models.WithdrawalRequest) (models.WithdrawalResponse, error) {
	return models.WithdrawalResponse{
		TransactionID: "test",
	}, nil
}
