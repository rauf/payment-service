package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rauf/payment-service/internal/backoff"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/protocol"
	"github.com/rauf/payment-service/internal/serde"
)

// GatewayA is the gateway for service A. It uses HTTP protocol with JSON serde
type GatewayA struct {
	baseGateway[gatewayARequest, gatewayAResponse]
}

func NewGatewayA(name, method, address string, httpClient *http.Client, retryConfig backoff.RetryConfig) *GatewayA {
	return &GatewayA{
		baseGateway: newBaseGateway[gatewayARequest, gatewayAResponse](
			name,
			serde.NewJSONSerde(),
			protocol.NewHTTPConnectionMock(httpClient, method, address, "json"),
			retryConfig,
		),
	}
}

func (g *GatewayA) Transact(ctx context.Context, transaction models.TransactionRequest) (models.TransactionResponse, error) {
	req := gatewayARequest{
		Amount:   transaction.Amount,
		Currency: transaction.Currency,
	}
	res, err := g.sendWithRetry(ctx, req)
	if err != nil {
		return models.TransactionResponse{}, fmt.Errorf("error sending transaction request: %w", err)
	}

	return models.TransactionResponse{
		RefID:     res.RefID,
		Status:    res.Status,
		CreatedAt: res.CreatedAt,
	}, nil
}
