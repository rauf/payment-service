package gateway

import (
	"context"
	"fmt"

	"github.com/rauf/payment-service/internal/models"
)

var (
	ErrGatewayUnavailable = fmt.Errorf("gateway unavailable")
)

type PaymentGateway interface {
	Name() string
	Transact(context.Context, models.TransactionRequest) (models.TransactionResponse, error)
}
