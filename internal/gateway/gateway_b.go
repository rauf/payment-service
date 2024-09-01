package gateway

import (
	"context"

	"github.com/rauf/payment-service/internal/format"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/protocol"
)

type GatewayB struct {
	baseGateway
}

func NewGatewayISO8583(address string) *GatewayB {
	return &GatewayB{
		baseGateway: baseGateway{
			dataFormat:      format.NewISO8583Protocol(),
			protocolHandler: protocol.NewTCPConnection(address),
		},
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
