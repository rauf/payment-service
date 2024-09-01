package handlers

import (
	"time"

	"github.com/rauf/payment-service/internal/validation"
)

type (
	depositApiRequest struct {
		Amount           float64           `json:"amount" validate:"required,gt=0"`
		Currency         string            `json:"currency" validate:"required,len=3"`
		PaymentMethod    string            `json:"payment_method" validate:"required"`
		Description      string            `json:"description,omitempty"`
		CustomerID       string            `json:"customer_id" validate:"required,uuid"`
		PreferredGateway string            `json:"preferred_gateway"`
		Metadata         map[string]string `json:"metadata,omitempty"`
	}
	depositApiResponse struct {
		TransactionID string    `json:"transaction_id"`
		Status        string    `json:"status"`
		CreatedAt     time.Time `json:"created_at"`
	}
	withdrawalApiRequest struct {
		Amount           float64           `json:"amount" validate:"required,gt=0"`
		Currency         string            `json:"currency" validate:"required,len=3"`
		PaymentMethod    string            `json:"payment_method" validate:"required"`
		Description      string            `json:"description,omitempty"`
		CustomerID       string            `json:"customer_id" validate:"required,uuid"`
		PreferredGateway string            `json:"preferred_gateway"`
		Metadata         map[string]string `json:"metadata,omitempty"`
	}
	withdrawalApiResponse struct {
		TransactionID string    `json:"transaction_id"`
		Status        string    `json:"status"`
		CreatedAt     time.Time `json:"created_at"`
	}
	updateStatusApiRequest struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"`
	}
	updateStatusApiResponse struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"`
	}
)

func (d *depositApiRequest) validate() validation.Errors {
	var errors validation.Errors
	if d.Amount <= 0 {
		errors.Add("amount", "must be greater than 0")
	}
	if len(d.Currency) != 3 {
		errors.Add("currency", "must be 3 characters long")
	}
	if d.PaymentMethod == "" {
		errors.Add("payment_method", "cannot be empty")
	}
	if d.CustomerID == "" {
		errors.Add("customer_id", "cannot be empty")
	}
	return errors
}

func (w *withdrawalApiRequest) validate() validation.Errors {
	var errors validation.Errors
	if w.Amount <= 0 {
		errors.Add("amount", "must be greater than 0")
	}
	if len(w.Currency) != 3 {
		errors.Add("currency", "must be 3 characters long")
	}
	if w.PaymentMethod == "" {
		errors.Add("payment_method", "cannot be empty")
	}
	if w.CustomerID == "" {
		errors.Add("customer_id", "cannot be empty")
	}
	return errors
}

func (c *updateStatusApiRequest) validate() validation.Errors {
	var errors validation.Errors
	if c.TransactionID == "" {
		errors.Add("transaction_id", "cannot be empty")
	}
	if c.Status == "" {
		errors.Add("status", "cannot be empty")
	}
	return errors
}
