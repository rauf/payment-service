package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/service"
	"github.com/rauf/payment-service/internal/utils/jsonutil"
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
	slog.InfoContext(r.Context(), "Deposit request received", "method", r.Method, "url", r.URL.Path)

	var apiRequest depositApiRequest
	if err := jsonutil.ReadJSON(r, &apiRequest); err != nil {
		return NewResponse(http.StatusBadRequest, "failed to decode request", nil, err)
	}
	if validationErrs := apiRequest.validate(); !validationErrs.IsValid() {
		return NewResponse(http.StatusBadRequest, "failed to validate request", validationErrs, &validationErrs)
	}

	req := models.DepositRequest{
		Amount:           apiRequest.Amount,
		Currency:         apiRequest.Currency,
		PaymentMethod:    apiRequest.PaymentMethod,
		Description:      apiRequest.Description,
		CustomerID:       apiRequest.CustomerID,
		PreferredGateway: apiRequest.PreferredGateway,
		Metadata:         apiRequest.Metadata,
	}

	res, err := h.paymentService.Deposit(r.Context(), req)
	if err != nil {
		if errors.Is(err, gateway.ErrGatewayUnavailable) {
			return NewResponse(http.StatusServiceUnavailable, "all payment gateways are currently unavailable", nil, err)
		}
		return NewResponse(http.StatusInternalServerError, "failed to process deposit", nil, err)
	}
	apiResponse := depositApiResponse{
		TransactionID: res.TransactionID,
		Status:        res.Status,
		CreatedAt:     res.CreatedAt,
	}
	return NewResponse(http.StatusOK, "deposit successful", apiResponse, nil)
}

func (h *PaymentHandler) HandleWithdrawal(w http.ResponseWriter, r *http.Request) error {
	slog.InfoContext(r.Context(), "Withdrawal request received", "method", r.Method, "url", r.URL.Path)

	var apiRequest withdrawalApiRequest
	if err := jsonutil.ReadJSON(r, &apiRequest); err != nil {
		return NewResponse(http.StatusBadRequest, "failed to decode request", nil, err)
	}
	if validationErrs := apiRequest.validate(); !validationErrs.IsValid() {
		return NewResponse(http.StatusBadRequest, "failed to validate request", validationErrs, &validationErrs)
	}

	req := models.WithdrawalRequest{
		Amount:           apiRequest.Amount,
		Currency:         apiRequest.Currency,
		PaymentMethod:    apiRequest.PaymentMethod,
		Description:      apiRequest.Description,
		CustomerID:       apiRequest.CustomerID,
		PreferredGateway: apiRequest.PreferredGateway,
		Metadata:         apiRequest.Metadata,
	}

	res, err := h.paymentService.Withdraw(r.Context(), req)
	if err != nil {
		if errors.Is(err, gateway.ErrGatewayUnavailable) {
			return NewResponse(http.StatusServiceUnavailable, "all payment gateways are currently unavailable", nil, err)
		}
		return NewResponse(http.StatusInternalServerError, "failed to process withdrawal", nil, err)
	}
	apiResponse := withdrawalApiResponse{
		TransactionID: res.TransactionID,
		Status:        res.Status,
		CreatedAt:     res.CreatedAt,
	}

	return NewResponse(http.StatusOK, "withdrawal successful", apiResponse, nil)
}

func (h *PaymentHandler) HandleUpdateStatus(w http.ResponseWriter, r *http.Request) error {
	slog.InfoContext(r.Context(), "Callback request received", "method", r.Method, "url", r.URL.Path)

	var apiRequest updateStatusApiRequest
	if err := jsonutil.ReadJSON(r, &apiRequest); err != nil {
		return NewResponse(http.StatusBadRequest, "failed to decode request", nil, err)
	}
	if validationErrs := apiRequest.validate(); !validationErrs.IsValid() {
		return NewResponse(http.StatusBadRequest, "failed to validate request", validationErrs, &validationErrs)
	}

	req := models.UpdateStatusRequest{
		TransactionID: apiRequest.TransactionID,
		Status:        apiRequest.Status,
	}

	err := h.paymentService.UpdateStatus(r.Context(), req)
	if err != nil {
		return NewResponse(http.StatusInternalServerError, "failed to process update status", nil, err)
	}
	apiResponse := updateStatusApiResponse{
		TransactionID: req.TransactionID,
		Status:        req.Status,
	}

	return NewResponse(http.StatusOK, "status updated successfully", apiResponse, nil)
}
