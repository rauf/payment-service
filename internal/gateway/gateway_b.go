package gateway

import (
	"context"
	"time"

	"github.com/rauf/payment-service/internal/backoff"
	"github.com/rauf/payment-service/internal/config"
	"github.com/rauf/payment-service/internal/format"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/protocol"
)

type GatewayB struct {
	baseGateway
}

func NewGatewayISO8583(address string) *GatewayB {
	return &GatewayB{
		baseGateway: newBaseGateway(
			"Gateway-B",
			format.NewISO8583Protocol(),
			protocol.NewTCPConnection(address),
			config.RetryConfig{
				MaxRetries: 3,
				Backoff:    backoff.NewExponentialBackoff(1*time.Second, 1.2, 2*time.Second),
			},
		),
	}
}

func (g *GatewayB) Deposit(ctx context.Context, deposit models.DepositRequest) (models.DepositResponse, error) {
	return models.DepositResponse{
		TransactionID: "test",
	}, nil
}

func (g *GatewayB) Withdraw(ctx context.Context, withdrawal models.WithdrawalRequest) (models.WithdrawalResponse, error) {
	return models.WithdrawalResponse{
		TransactionID: "test",
	}, nil
}
