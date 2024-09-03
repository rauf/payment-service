package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CreateTransaction(ctx context.Context, req models.TransactionRequest) (models.TransactionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(models.TransactionResponse), args.Error(1)
}

func (m *MockPaymentService) UpdateStatus(ctx context.Context, req models.UpdateStatusRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestHandleTransaction(t *testing.T) {
	tests := []struct {
		name               string
		input              transactionApiRequest
		mockResponse       models.TransactionResponse
		mockError          error
		callTransactMethod bool
		expectedStatus     int
		expectedBody       string
	}{
		{
			name: "Successful transaction",
			input: transactionApiRequest{
				Amount:           100.0,
				Type:             "deposit",
				Currency:         "USD",
				PaymentMethod:    "card",
				CustomerID:       "cust123",
				PreferredGateway: "stripe",
			},
			mockResponse: models.TransactionResponse{
				RefID:   "ref123",
				Status:  "pending",
				Gateway: "stripe",
			},
			mockError:          nil,
			expectedStatus:     http.StatusOK,
			callTransactMethod: true,
			expectedBody:       `{"code":200,"message":"transaction sent to gateway successfully","data":{"ref_id":"ref123","status":"pending","created_at":"0001-01-01T00:00:00Z","gateway":"stripe"}}`,
		},
		{
			name: "Invalid request",
			input: transactionApiRequest{
				Amount: -100.0,
			},
			expectedStatus:     http.StatusBadRequest,
			callTransactMethod: false,
			expectedBody:       `{"code":400,"message":"failed to validate request","data":{"errors":[{"field":"amount","message":"must be greater than 0"},{"field":"currency","message":"must be 3 characters long"},{"field":"payment_method","message":"cannot be empty"},{"field":"customer_id","message":"cannot be empty"},{"field":"type","message":"cannot be empty"}]}}`,
		},
		{
			name: "Gateway unavailable",
			input: transactionApiRequest{
				Amount:           100.0,
				Type:             "deposit",
				Currency:         "USD",
				PaymentMethod:    "card",
				CustomerID:       "cust123",
				PreferredGateway: "stripe",
			},
			mockError:          gateway.ErrGatewayUnavailable,
			expectedStatus:     http.StatusServiceUnavailable,
			callTransactMethod: true,
			expectedBody:       `{"code":503,"message":"all payment gateways are currently unavailable"}`,
		},
	}

	for _, tt := range tests {
		mockService := new(MockPaymentService)
		handler := NewPaymentHandler(mockService)

		t.Run(tt.name, func(t *testing.T) {
			if tt.callTransactMethod {
				mockService.On("CreateTransaction", mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)
			}

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/transaction", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			res := handler.HandleTransaction(rr, req)
			writeResponse(rr, req, res)

			assert.Equal(t, tt.expectedStatus, res.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandleUpdateStatus(t *testing.T) {

	tests := []struct {
		name             string
		input            updateStatusApiRequest
		mockError        error
		callUpdateMethod bool
		expectedStatus   int
		expectedBody     string
	}{
		{
			name: "Successful status update",
			input: updateStatusApiRequest{
				Gateway: "stripe",
				RefID:   "ref123",
				Status:  "success",
			},
			mockError:        nil,
			callUpdateMethod: true,
			expectedStatus:   http.StatusOK,
			expectedBody:     `{"code":200,"message":"status updated successfully"}`,
		},
		{
			name: "Invalid request",
			input: updateStatusApiRequest{
				Gateway: "",
				RefID:   "",
				Status:  "invalid",
			},
			mockError:        nil,
			callUpdateMethod: false,
			expectedStatus:   http.StatusBadRequest,
			expectedBody:     `{"code":400,"message":"failed to validate request","data":{"errors":[{"field":"ref_id","message":"cannot be empty"},{"field":"gateway","message":"cannot be empty"},{"field":"status","message":"not valid transaction status"}]}}`,
		},
		{
			name: "Transaction not found",
			input: updateStatusApiRequest{
				Gateway: "stripe",
				RefID:   "ref1234",
				Status:  "success",
			},
			mockError:        service.ErrTransactionNotFound,
			callUpdateMethod: true,
			expectedStatus:   http.StatusNotFound,
			expectedBody:     `{"code":404,"message":"transaction not found"}`,
		},
	}

	for _, tt := range tests {
		mockService := new(MockPaymentService)
		handler := NewPaymentHandler(mockService)

		t.Run(tt.name, func(t *testing.T) {
			if tt.callUpdateMethod {
				mockService.On("UpdateStatus", mock.Anything, mock.Anything).Return(tt.mockError)
			}
			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/update-status", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			res := handler.HandleUpdateStatus(rr, req)
			writeResponse(rr, req, res)

			assert.Equal(t, tt.expectedStatus, res.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
