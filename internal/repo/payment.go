package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/utils/nullutil"
)

type PaymentRepo struct {
	queries *models.Queries
}

func NewPaymentRepo(queries *models.Queries) *PaymentRepo {
	return &PaymentRepo{
		queries: queries,
	}
}

func (r *PaymentRepo) CreateTransaction(ctx context.Context, transaction CreateTransaction) error {
	arg := models.CreateTransactionParams{
		Type:             models.TransactionType(strings.ToUpper(transaction.Type)),
		Amount:           fmt.Sprintf("%.2f", transaction.Amount),
		Currency:         transaction.Currency,
		PaymentMethod:    transaction.PaymentMethod,
		Description:      nullutil.NewNullString(transaction.Description),
		CustomerID:       transaction.CustomerID,
		Gateway:          transaction.Gateway,
		GatewayRefID:     transaction.GatewayRefID,
		Status:           models.TransactionStatusPENDING,
		PreferredGateway: nullutil.NewNullString(transaction.PreferredGateway),
		Metadata:         nullutil.NewNullRawMessage(transaction.Metadata),
	}

	return r.queries.CreateTransaction(ctx, arg)
}

func (r *PaymentRepo) GetTransactionByRefID(ctx context.Context, g GetTransactionByRefID) (models.Transaction, error) {
	return r.queries.GetTransactionByGatewayRefId(ctx, models.GetTransactionByGatewayRefIdParams{
		GatewayRefID: g.RefID,
		Gateway:      g.Gateway,
	})
}

func (r *PaymentRepo) UpdateTransactionStatus(ctx context.Context, update UpdateTransactionStatus) error {
	arg := models.UpdateTransactionStatusParams{
		Gateway:      update.Gateway,
		GatewayRefID: update.RefID,
		Status:       models.TransactionStatus(strings.ToUpper(update.Status)),
		UpdatedAt:    time.Now().UTC(),
	}

	return r.queries.UpdateTransactionStatus(ctx, arg)
}
