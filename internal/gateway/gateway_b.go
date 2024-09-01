package gateway

import (
	"context"
	"time"

	"github.com/rauf/payment-service/internal/backoff"
	"github.com/rauf/payment-service/internal/config"
	"github.com/rauf/payment-service/internal/format"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/protocol"
	"github.com/rauf/payment-service/internal/utils/randutil"
)

type GatewayB struct {
	baseGateway[gatewayBRequest, gatewayBResponse]
}

func NewGatewayISO8583(address string) *GatewayB {
	return &GatewayB{
		baseGateway: newBaseGateway[gatewayBRequest, gatewayBResponse](
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
		TransactionID: randutil.RandomString(10),
	}, nil
}

func (g *GatewayB) Withdraw(ctx context.Context, withdrawal models.WithdrawalRequest) (models.WithdrawalResponse, error) {
	return models.WithdrawalResponse{
		TransactionID: randutil.RandomString(10),
	}, nil
}
