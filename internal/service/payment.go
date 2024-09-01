package service

import (
	"context"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/router"
	"github.com/rauf/payment-service/internal/utils"
)

type PaymentService struct {
	router *router.Router
}

func NewPaymentService(router *router.Router) *PaymentService {
	return &PaymentService{
		router: router,
	}
}

func (s *PaymentService) Deposit(ctx context.Context, deposit models.DepositRequest) (models.DepositResponse, error) {
	response, err := s.router.SendMessage(deposit.PreferredGateway, func(g gateway.PaymentGateway) (any, error) {
		return g.Deposit(ctx, deposit)
	})
	if err != nil {
		return models.DepositResponse{}, err
	}
	return utils.Cast[models.DepositResponse](response)
}

func (s *PaymentService) Withdraw(ctx context.Context, withdrawal models.WithdrawalRequest) (models.WithdrawalResponse, error) {
	response, err := s.router.SendMessage(withdrawal.PreferredGateway, func(g gateway.PaymentGateway) (any, error) {
		return g.Withdraw(ctx, withdrawal)
	})
	if err != nil {
		return models.WithdrawalResponse{}, err
	}
	return utils.Cast[models.WithdrawalResponse](response)
}
