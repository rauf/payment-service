package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) HandleDeposit(w http.ResponseWriter, r *http.Request) error {
	resp, err := h.paymentService.Deposit(r.Context(), models.DepositRequest{
		Amount:           0,
		Currency:         "",
		PreferredGateway: "",
	})

	if err != nil {
		return err
	}
	jsonBody, err := json.Marshal(resp)

	if err != nil {
		return err
	}
	if _, err := w.Write(jsonBody); err != nil {
		return err
	}
	return nil
}
