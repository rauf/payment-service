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
	"github.com/rauf/payment-service/internal/utils/randutil"
)

type GatewayA struct {
	baseGateway[gatewayARequest, gatewayAResponse]
}

func NewGatewayA(address string) *GatewayA {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &GatewayA{
		baseGateway: newBaseGateway[gatewayARequest, gatewayAResponse](
			"Gateway-A",
			format.NewJSONProtocol(),
			protocol.NewHTTPConnectionMock(httpClient, http.MethodPost, address),
			config.RetryConfig{
				MaxRetries: 3,
				Backoff:    backoff.NewExponentialBackoff(1*time.Second, 1.2, 2*time.Second),
			},
		),
	}
}

func (g *GatewayA) Deposit(ctx context.Context, deposit models.DepositRequest) (models.DepositResponse, error) {
	req := gatewayARequest{
		Amount:   deposit.Amount,
		Currency: deposit.Currency,
	}
	res, err := g.SendWithRetry(ctx, req)
	if err != nil {
		return models.DepositResponse{}, fmt.Errorf("error sending deposit request: %w", err)
	}

	return models.DepositResponse{
		TransactionID: res.TransactionID,
		Status:        res.Status,
		CreatedAt:     res.CreatedAt,
	}, nil
}

func (g *GatewayA) Withdraw(ctx context.Context, withdrawal models.WithdrawalRequest) (models.WithdrawalResponse, error) {
	return models.WithdrawalResponse{
		TransactionID: randutil.RandomString(10),
	}, nil
}
