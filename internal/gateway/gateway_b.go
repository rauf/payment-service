package gateway

import (
	"context"
	"net/http"

	"github.com/rauf/payment-service/internal/backoff"
	"github.com/rauf/payment-service/internal/format"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/protocol"
	"github.com/rauf/payment-service/internal/utils/randutil"
)

type GatewayB struct {
	baseGateway[gatewayBRequest, gatewayBResponse]
}

func NewGatewayB(name, method, address string, httpClient *http.Client, retryConfig backoff.RetryConfig) *GatewayB {
	return &GatewayB{
		baseGateway: newBaseGateway[gatewayBRequest, gatewayBResponse](
			name,
			format.NewISO8583Protocol(),
			protocol.NewHTTPConnectionMock(httpClient, method, address),
			retryConfig,
		),
	}
}

func (g *GatewayB) Transact(ctx context.Context, transaction models.TransactionRequest) (models.TransactionResponse, error) {
	return models.TransactionResponse{
		RefID: randutil.RandomString(10),
	}, nil
}
