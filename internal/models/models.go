package models

import "time"

type DepositRequest struct {
	Amount           float64
	Currency         string
	PaymentMethod    string
	Description      string
	CustomerID       string
	PreferredGateway string

	Metadata map[string]string
}

type DepositResponse struct {
	TransactionID string
	Status        string
	CreatedAt     time.Time
}

type WithdrawalRequest struct {
	Amount           float64
	Currency         string
	PaymentMethod    string
	Description      string
	CustomerID       string
	PreferredGateway string

	Metadata map[string]string
}

type WithdrawalResponse struct {
	TransactionID string
	Status        string
	CreatedAt     time.Time
}

type UpdateStatusRequest struct {
	TransactionID string
	Status        string
}

type UpdateStatusResponse struct {
	TransactionID string
	Status        string
}
