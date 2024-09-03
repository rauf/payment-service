package handlers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/rauf/payment-service/internal/validation"
)

var (
	allowedTransactionTypes = map[string]struct{}{
		"deposit":    {},
		"withdrawal": {},
	}
	allowedTransactionStatuses = map[string]struct{}{
		"pending": {},
		"success": {},
		"failed":  {},
	}
)

type (
	validate interface {
		validate() validation.Errors
	}
	transactionApiRequest struct {
		Amount           float64         `json:"amount"`
		Type             string          `json:"type"`
		Currency         string          `json:"currency"`
		PaymentMethod    string          `json:"payment_method"`
		Description      string          `json:"description,omitempty"`
		CustomerID       string          `json:"customer_id"`
		PreferredGateway string          `json:"preferred_gateway"`
		Metadata         json.RawMessage `json:"metadata,omitempty"`
	}
	transactionApiResponse struct {
		RefID     string    `json:"ref_id"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		Gateway   string    `json:"gateway"`
	}
	updateStatusApiRequest struct {
		Gateway string `json:"gateway"`
		RefID   string `json:"ref_id"`
		Status  string `json:"status"`
	}
	gatewayACallbackRequest struct {
		RefID     string    `json:"ref_id"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}
	gatewayBCallbackRequest struct {
		RefID     string    `xml:"ref_id"`
		Status    string    `xml:"status"`
		CreatedAt time.Time `xml:"created_at"`
	}
)

func (d *transactionApiRequest) validate() validation.Errors {
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
	if d.Type == "" {
		errors.Add("type", "cannot be empty")
	} else if _, ok := allowedTransactionTypes[strings.ToLower(d.Type)]; !ok {
		errors.Add("type", "not valid transaction type")
	}
	return errors
}

func (c *updateStatusApiRequest) validate() validation.Errors {
	var errors validation.Errors
	if c.RefID == "" {
		errors.Add("ref_id", "cannot be empty")
	}
	if c.Gateway == "" {
		errors.Add("gateway", "cannot be empty")
	}
	if c.Status == "" {
		errors.Add("status", "cannot be empty")
	} else if _, ok := allowedTransactionStatuses[strings.ToLower(c.Status)]; !ok {
		errors.Add("status", "not valid transaction status")
	}
	return errors
}

func (r gatewayACallbackRequest) validate() validation.Errors {
	var errors validation.Errors
	if r.RefID == "" {
		errors.Add("ref_id", "cannot be empty")
	}
	if r.Status == "" {
		errors.Add("status", "cannot be empty")
	} else if _, ok := allowedTransactionStatuses[strings.ToLower(r.Status)]; !ok {
		errors.Add("status", "not valid transaction status")
	}
	return errors
}

func (r gatewayBCallbackRequest) validate() validation.Errors {
	var errors validation.Errors
	if r.RefID == "" {
		errors.Add("ref_id", "cannot be empty")
	}
	if r.Status == "" {
		errors.Add("status", "cannot be empty")
	} else if _, ok := allowedTransactionStatuses[strings.ToLower(r.Status)]; !ok {
		errors.Add("status", "not valid transaction status")
	}
	return errors
}
