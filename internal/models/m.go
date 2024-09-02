package models

import (
	"encoding/json"
	"time"
)

type TransactionRequest struct {
	Type             string
	Amount           float64
	Currency         string
	PaymentMethod    string
	Description      string
	CustomerID       string
	PreferredGateway string
	Metadata         json.RawMessage
}

type TransactionResponse struct {
	Gateway   string
	RefID     string
	Status    string
	CreatedAt time.Time
}

type UpdateStatusRequest struct {
	Gateway string
	RefID   string
	Status  string
}

type UpdateStatusResponse struct {
	RefID  string
	Status string
}
