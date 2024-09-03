package gateway

import "time"

type gatewayARequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type gatewayAResponse struct {
	RefID     string    `json:"ref_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type gatewayBRequest struct {
	Amount   float64 `xml:"amount"`
	Currency string  `xml:"currency"`
}

type gatewayBResponse struct {
	RefID     string    `xml:"ref_id"`
	Status    string    `xml:"status"`
	CreatedAt time.Time `xml:"created_at"`
}
