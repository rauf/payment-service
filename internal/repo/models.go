package repo

import "github.com/rauf/payment-service/internal/models"

type CreateTransaction struct {
	models.TransactionRequest
	Gateway      string
	GatewayRefID string
}

type GetTransactionByRefID struct {
	Gateway string
	RefID   string
}

type UpdateTransactionStatus struct {
	Gateway string
	RefID   string
	Status  string
}
