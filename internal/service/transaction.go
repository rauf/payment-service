package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/repo"
	"github.com/rauf/payment-service/internal/router"
)

var ErrTransactionNotFound = errors.New("transaction not found")

// PaymentService is a service that handles payment transactions
type PaymentService struct {
	router      *router.Router
	paymentRepo *repo.PaymentRepo
}

func NewPaymentService(router *router.Router, paymentRepo *repo.PaymentRepo) *PaymentService {
	return &PaymentService{
		router:      router,
		paymentRepo: paymentRepo,
	}
}

func (s *PaymentService) CreateTransaction(ctx context.Context, transaction models.TransactionRequest) (models.TransactionResponse, error) {
	response, err := s.router.SendMessage(ctx, transaction.PreferredGateway, func(g gateway.PaymentGateway) (models.TransactionResponse, error) {
		return g.Transact(ctx, transaction)
	})

	slog.InfoContext(ctx, "Received response from gateway", "response", response, "error", err)

	if errors.Is(err, gateway.ErrGatewayUnavailable) {
		return models.TransactionResponse{}, fmt.Errorf("all payment gateways are currently unavailable: %w", err)
	}
	if err != nil {
		return models.TransactionResponse{}, fmt.Errorf("transaction failed: %w", err)
	}
	err = s.paymentRepo.CreateTransaction(ctx, repo.CreateTransaction{
		TransactionRequest: transaction,
		Gateway:            response.Gateway,
		GatewayRefID:       response.Data.RefID,
	})
	if err != nil {
		return models.TransactionResponse{}, fmt.Errorf("failed to save transaction: %w", err)
	}
	return models.TransactionResponse{
		Gateway:   response.Gateway,
		RefID:     response.Data.RefID,
		Status:    response.Data.Status,
		CreatedAt: response.Data.CreatedAt,
	}, nil
}

func (s *PaymentService) UpdateStatus(ctx context.Context, req models.UpdateStatusRequest) error {
	slog.InfoContext(ctx, "Updating transaction status", "ref_id", req.RefID, "status", req.Status)

	_, err := s.paymentRepo.GetTransactionByRefID(ctx, repo.GetTransactionByRefID{
		Gateway: req.Gateway,
		RefID:   req.RefID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: transaction with ref_id %s not found", ErrTransactionNotFound, req.RefID)
		}
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = s.paymentRepo.UpdateTransactionStatus(ctx, repo.UpdateTransactionStatus{
		Gateway: req.Gateway,
		RefID:   req.RefID,
		Status:  req.Status,
	})
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	return nil
}
