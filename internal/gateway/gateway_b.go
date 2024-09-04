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

// GatewayB is the gateway for service B. It uses HTTP protocol with XML serde
type GatewayB struct {
	baseGateway[gatewayBRequest, gatewayBResponse]
}

func NewGatewayB(name, method, address string, httpClient *http.Client, retryConfig backoff.RetryConfig) *GatewayB {
	return &GatewayB{
		baseGateway: newBaseGateway[gatewayBRequest, gatewayBResponse](
			name,
			serde.NewXMLSerde(),
			protocol.NewHTTPConnectionMock(httpClient, method, address, "xml"),
			retryConfig,
		),
	}
}

func (g *GatewayB) Transact(ctx context.Context, transaction models.TransactionRequest) (models.TransactionResponse, error) {
	req := gatewayBRequest{
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
