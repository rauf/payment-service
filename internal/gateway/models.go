package gateway

import "time"

type gatewayARequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type gatewayAResponse struct {
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type gatewayBRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type gatewayBResponse struct {
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
