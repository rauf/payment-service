package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/service"
	"github.com/rauf/payment-service/internal/utils/jsonutil"
)

type PaymentHandler struct {
	paymentService paymentService
}

// interface on consumer side
type paymentService interface {
	CreateTransaction(ctx context.Context, req models.TransactionRequest) (models.TransactionResponse, error)
	UpdateStatus(ctx context.Context, req models.UpdateStatusRequest) error
}

func NewPaymentHandler(paymentService paymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) HandleTransaction(w http.ResponseWriter, r *http.Request) Response {
	slog.InfoContext(r.Context(), "Transaction request received", "method", r.Method, "url", r.URL.Path)

	var apiRequest transactionApiRequest
	if err := jsonutil.ReadJSON(r, &apiRequest); err != nil {
		return NewResponse(http.StatusBadRequest, "failed to decode request", nil, err)
	}
	if validationErrs := apiRequest.validate(); !validationErrs.IsValid() {
		return NewResponse(http.StatusBadRequest, "failed to validate request", validationErrs, &validationErrs)
	}

	req := models.TransactionRequest{
		Type:             apiRequest.Type,
		Amount:           apiRequest.Amount,
		Currency:         apiRequest.Currency,
		PaymentMethod:    apiRequest.PaymentMethod,
		Description:      apiRequest.Description,
		CustomerID:       apiRequest.CustomerID,
		PreferredGateway: apiRequest.PreferredGateway,
		Metadata:         apiRequest.Metadata,
	}

	res, err := h.paymentService.CreateTransaction(r.Context(), req)
	if err != nil {
		if errors.Is(err, gateway.ErrGatewayUnavailable) {
			return NewResponse(http.StatusServiceUnavailable, "all payment gateways are currently unavailable", nil, err)
		}
		return NewResponse(http.StatusInternalServerError, "failed to process transaction", nil, err)
	}
	apiResponse := transactionApiResponse{
		RefID:     res.RefID,
		Status:    res.Status,
		CreatedAt: res.CreatedAt,
		Gateway:   res.Gateway,
	}
	return NewResponse(http.StatusOK, "transaction sent to gateway successfully", apiResponse, nil)
}

func (h *PaymentHandler) HandleUpdateStatus(w http.ResponseWriter, r *http.Request) Response {
	slog.InfoContext(r.Context(), "Callback request received", "method", r.Method, "url", r.URL.Path)

	var apiRequest updateStatusApiRequest
	if err := jsonutil.ReadJSON(r, &apiRequest); err != nil {
		return NewResponse(http.StatusBadRequest, "failed to decode request", nil, err)
	}
	if validationErrs := apiRequest.validate(); !validationErrs.IsValid() {
		return NewResponse(http.StatusBadRequest, "failed to validate request", validationErrs, &validationErrs)
	}

	req := models.UpdateStatusRequest{
		Gateway: apiRequest.Gateway,
		RefID:   apiRequest.RefID,
		Status:  apiRequest.Status,
	}

	err := h.paymentService.UpdateStatus(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrTransactionNotFound) {
			return NewResponse(http.StatusNotFound, "transaction not found", nil, err)
		}
		return NewResponse(http.StatusInternalServerError, "failed to process update status", nil, err)
	}
	return NewResponse(http.StatusOK, "status updated successfully", nil, nil)
}
