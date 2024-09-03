package gateway

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rauf/payment-service/internal/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDataFormat struct {
	mock.Mock
}

func (m *mockDataFormat) Marshal(v interface{}) ([]byte, error) {
	args := m.Called(v)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockDataFormat) Unmarshal(data []byte, v interface{}) error {
	args := m.Called(data, v)
	return args.Error(0)
}

type mockProtocolHandler struct {
	mock.Mock
}

func (m *mockProtocolHandler) Send(ctx context.Context, data []byte) ([]byte, error) {
	args := m.Called(ctx, data)
	return args.Get(0).([]byte), args.Error(1)
}

func TestNewBaseGateway(t *testing.T) {
	name := "test_gateway"
	dataFormat := &mockDataFormat{}
	protocolHandler := &mockProtocolHandler{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}

	bg := newBaseGateway[string, int](name, dataFormat, protocolHandler, retryConfig)

	assert.Equal(t, name, bg.name)
	assert.Equal(t, dataFormat, bg.dataFormat)
	assert.Equal(t, protocolHandler, bg.protocolHandler)
	assert.Equal(t, retryConfig, bg.retryConfig)
}

func TestBaseGateway_Send_Success(t *testing.T) {
	dataFormat := &mockDataFormat{}
	protocolHandler := &mockProtocolHandler{}
	bg := newBaseGateway[string, int]("test", dataFormat, protocolHandler, backoff.RetryConfig{})

	dataFormat.On("Marshal", "test_data").Return([]byte("encoded_data"), nil)
	protocolHandler.On("Send", mock.Anything, []byte("encoded_data")).Return([]byte("response_data"), nil)
	dataFormat.On("Unmarshal", []byte("response_data"), mock.AnythingOfType("*int")).Return(nil).Run(func(args mock.Arguments) {
		*(args.Get(1).(*int)) = 42
	})

	result, err := bg.send(context.Background(), "test_data")

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
	dataFormat.AssertExpectations(t)
	protocolHandler.AssertExpectations(t)
}

func TestBaseGateway_Send_MarshalError(t *testing.T) {
	dataFormat := &mockDataFormat{}
	protocolHandler := &mockProtocolHandler{}
	bg := newBaseGateway[string, int]("test", dataFormat, protocolHandler, backoff.RetryConfig{})

	dataFormat.On("Marshal", "test_data").Return([]byte{}, errors.New("marshal error"))

	_, err := bg.send(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error marshaling data")
	dataFormat.AssertExpectations(t)
}

func TestBaseGateway_Send_SendError(t *testing.T) {
	dataFormat := &mockDataFormat{}
	protocolHandler := &mockProtocolHandler{}
	bg := newBaseGateway[string, int]("test", dataFormat, protocolHandler, backoff.RetryConfig{})

	dataFormat.On("Marshal", "test_data").Return([]byte("encoded_data"), nil)
	protocolHandler.On("Send", mock.Anything, []byte("encoded_data")).Return([]byte{}, errors.New("send error"))

	_, err := bg.send(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error sending data")
	dataFormat.AssertExpectations(t)
	protocolHandler.AssertExpectations(t)
}

func TestBaseGateway_SendWithRetry_Success(t *testing.T) {
	dataFormat := &mockDataFormat{}
	protocolHandler := &mockProtocolHandler{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", dataFormat, protocolHandler, retryConfig)

	dataFormat.On("Marshal", "test_data").Return([]byte("encoded_data"), nil)
	protocolHandler.On("Send", mock.Anything, []byte("encoded_data")).Return([]byte("response_data"), nil)
	dataFormat.On("Unmarshal", []byte("response_data"), mock.AnythingOfType("*int")).Return(nil).Run(func(args mock.Arguments) {
		*(args.Get(1).(*int)) = 42
	})

	result, err := bg.sendWithRetry(context.Background(), "test_data")

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
	dataFormat.AssertExpectations(t)
	protocolHandler.AssertExpectations(t)
}

func TestBaseGateway_SendWithRetry_EventualSuccess(t *testing.T) {
	dataFormat := &mockDataFormat{}
	protocolHandler := &mockProtocolHandler{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", dataFormat, protocolHandler, retryConfig)

	dataFormat.On("Marshal", "test_data").Return([]byte("encoded_data"), nil)
	protocolHandler.On("Send", mock.Anything, []byte("encoded_data")).Return([]byte{}, ErrGatewayUnavailable).Once()
	protocolHandler.On("Send", mock.Anything, []byte("encoded_data")).Return([]byte("response_data"), nil).Once()
	dataFormat.On("Unmarshal", []byte("response_data"), mock.AnythingOfType("*int")).Return(nil).Run(func(args mock.Arguments) {
		*(args.Get(1).(*int)) = 42
	})

	result, err := bg.sendWithRetry(context.Background(), "test_data")

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
	dataFormat.AssertExpectations(t)
	protocolHandler.AssertExpectations(t)
}

func TestBaseGateway_SendWithRetry_MaxRetriesReached(t *testing.T) {
	dataFormat := &mockDataFormat{}
	protocolHandler := &mockProtocolHandler{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 2,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", dataFormat, protocolHandler, retryConfig)

	dataFormat.On("Marshal", "test_data").Return([]byte("encoded_data"), nil)
	protocolHandler.On("Send", mock.Anything, []byte("encoded_data")).Return([]byte{}, ErrGatewayUnavailable)

	_, err := bg.sendWithRetry(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max retries reached")
	dataFormat.AssertExpectations(t)
	protocolHandler.AssertNumberOfCalls(t, "Send", 3) // Initial attempt + 2 retries
}

func TestBaseGateway_Name(t *testing.T) {
	bg := newBaseGateway[string, int]("test_gateway", nil, nil, backoff.RetryConfig{})
	assert.Equal(t, "test_gateway", bg.Name())

	bgNoName := newBaseGateway[string, int]("", nil, nil, backoff.RetryConfig{})
	assert.Equal(t, "unnamed gateway", bgNoName.Name())
}
