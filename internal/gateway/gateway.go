package gateway

import (
	"context"
	"fmt"

	"github.com/rauf/payment-service/internal/models"
)

var (
	// ErrGatewayUnavailable is an error that is returned when the gateway is unavailable.
	ErrGatewayUnavailable = fmt.Errorf("gateway unavailable")
)

// PaymentGateway is an interface that defines the methods that a payment gateway should implement.
type PaymentGateway interface {
	Name() string
	Transact(context.Context, models.TransactionRequest) (models.TransactionResponse, error)
}
