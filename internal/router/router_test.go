package router

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/registry"
	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockGateway struct {
	mock.Mock
	name string
}

func (m *mockGateway) Transact(ctx context.Context, request models.TransactionRequest) (models.TransactionResponse, error) {
	return m.Called().Get(0).(models.TransactionResponse), m.Called().Error(1)
}

func (m *mockGateway) Name() string {
	return m.name
}

func TestRouter_SendMessage(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		preferredGateway string
		gateways         []gateway.PaymentGateway
		operation        func(gateway.PaymentGateway) (models.TransactionResponse, error)
		expectedResponse Response
		expectedError    string
	}{
		{
			name:             "Success with preferred gateway",
			preferredGateway: "gateway1",
			gateways: []gateway.PaymentGateway{
				&mockGateway{name: "gateway1"},
				&mockGateway{name: "gateway2"},
			},
			operation: func(g gateway.PaymentGateway) (models.TransactionResponse, error) {
				return models.TransactionResponse{RefID: "123"}, nil
			},
			expectedResponse: Response{
				Gateway: "gateway1",
				Data:    models.TransactionResponse{RefID: "123"},
			},
		},
		{
			name:             "Fallback to second gateway",
			preferredGateway: "gateway1",
			gateways: []gateway.PaymentGateway{
				&mockGateway{name: "gateway1"},
				&mockGateway{name: "gateway2"},
			},
			operation: func(g gateway.PaymentGateway) (models.TransactionResponse, error) {
				if g.Name() == "gateway1" {
					return models.TransactionResponse{}, gateway.ErrGatewayUnavailable
				}
				return models.TransactionResponse{RefID: "456"}, nil
			},
			expectedResponse: Response{
				Gateway: "gateway2",
				Data:    models.TransactionResponse{RefID: "456"},
			},
		},
		{
			name:             "All gateways fail",
			preferredGateway: "gateway1",
			gateways: []gateway.PaymentGateway{
				&mockGateway{name: "gateway1"},
				&mockGateway{name: "gateway2"},
			},
			operation: func(g gateway.PaymentGateway) (models.TransactionResponse, error) {
				return models.TransactionResponse{}, gateway.ErrGatewayUnavailable
			},
			expectedError: "all gateways failed: gateway unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := registry.NewRegistry[gateway.PaymentGateway]()
			for _, g := range tt.gateways {
				_ = reg.Register(g.Name(), g)
			}

			r := NewRouter(reg, gobreaker.Settings{})

			response, err := r.SendMessage(ctx, tt.preferredGateway, tt.operation)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}

func TestRouter_fallback(t *testing.T) {
	reg := registry.NewRegistry[gateway.PaymentGateway]()

	require.NoError(t, reg.Register("GatewayA", &mockGateway{name: "GatewayA"}))
	require.NoError(t, reg.Register("GatewayB", &mockGateway{name: "GatewayB"}))

	// Create router
	router := NewRouter(reg, gobreaker.Settings{
		Name:    "TestCircuitBreaker",
		Timeout: 5 * time.Second,
	})

	testCases := []struct {
		name             string
		preferredGateway string
		operation        func(gateway.PaymentGateway) (models.TransactionResponse, error)
		expectedGateway  string
	}{
		{
			name:             "Preferred Gateway A",
			preferredGateway: "GatewayA",
			operation: func(paymentGateway gateway.PaymentGateway) (models.TransactionResponse, error) {
				return models.TransactionResponse{RefID: "123"}, nil
			},
			expectedGateway: "GatewayA",
		},
		{
			name:             "Preferred Gateway B",
			preferredGateway: "GatewayB",
			operation: func(paymentGateway gateway.PaymentGateway) (models.TransactionResponse, error) {
				return models.TransactionResponse{RefID: "123"}, nil
			},
			expectedGateway: "GatewayB",
		},
		{
			name:             "Fallback to Gateway B",
			preferredGateway: "NonExistentGateway",
			operation: func(paymentGateway gateway.PaymentGateway) (models.TransactionResponse, error) {
				return models.TransactionResponse{RefID: "123"}, nil
			},
			expectedGateway: "GatewayA", // Assuming GatewayA is the first in the list
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			response, err := router.SendMessage(ctx, tc.preferredGateway, tc.operation)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedGateway, response.Gateway)
			assert.NotEmpty(t, response.Data.RefID)
		})
	}
}

func TestRouterWithCircuitBreaker(t *testing.T) {
	failAfter := 3

	reg := registry.NewRegistry[gateway.PaymentGateway]()
	require.NoError(t, reg.Register("MockGateway", &mockGateway{name: "MockGateway"}))

	// Create router with circuit breaker settings
	router := NewRouter(reg, gobreaker.Settings{
		Name:     "TestCircuitBreaker",
		Interval: 1 * time.Second,
		Timeout:  1 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > uint32(failAfter)
		},
	})

	ctx := context.Background()

	// Test successful requests
	for i := 0; i < failAfter; i++ {
		response, err := router.SendMessage(ctx, "MockGateway", func(g gateway.PaymentGateway) (models.TransactionResponse, error) {
			return models.TransactionResponse{Gateway: "MockGateway", RefID: "mock-ref-id"}, nil
		})
		require.NoError(t, err)
		assert.Equal(t, "MockGateway", response.Gateway)
		assert.Equal(t, "mock-ref-id", response.Data.RefID)
	}

	// Test circuit breaker opening
	for i := 0; i < 5; i++ {
		_, err := router.SendMessage(ctx, "MockGateway", func(g gateway.PaymentGateway) (models.TransactionResponse, error) {
			return models.TransactionResponse{}, fmt.Errorf("error")
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "all gateways failed")
	}

	// Wait for the circuit breaker to close
	time.Sleep(2 * time.Second)

	// Test circuit breaker closing and successful request
	response, err := router.SendMessage(ctx, "MockGateway", func(g gateway.PaymentGateway) (models.TransactionResponse, error) {
		return models.TransactionResponse{Gateway: "MockGateway", RefID: "mock-ref-id"}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, "MockGateway", response.Gateway)
	assert.Equal(t, "mock-ref-id", response.Data.RefID)
}
