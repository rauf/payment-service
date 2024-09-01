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
	Deposit(context.Context, models.DepositRequest) (models.DepositResponse, error)
	Withdraw(context.Context, models.WithdrawalRequest) (models.WithdrawalResponse, error)
}
