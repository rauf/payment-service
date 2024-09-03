package gateway

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/rauf/payment-service/internal/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type mockSerde struct {
	mock.Mock
}

func (m *mockSerde) Serialize(w io.Writer, data any) error {
	args := m.Called(w, data)
	return args.Error(0)
}

func (m *mockSerde) Deserialize(r io.Reader, v any) error {
	args := m.Called(r, v)
	return args.Error(0)
}

type mockProtocol struct {
	mock.Mock
}

func (m *mockProtocol) Send(ctx context.Context, data []byte) ([]byte, error) {
	args := m.Called(ctx, data)
	return args.Get(0).([]byte), args.Error(1)
}

func TestBaseGateway_SendWithRetry_ContextCancelled(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, retryConfig)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(nil)
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte{}, errors.New("some error"))

	_, err := bg.sendWithRetry(ctx, "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context cancelled")
	mockSerde.AssertExpectations(t)
	mockProto.AssertExpectations(t)
}

func TestBaseGateway_Send_NilSerde(t *testing.T) {
	bg := newBaseGateway[string, int]("test", nil, &mockProtocol{}, backoff.RetryConfig{})

	_, err := bg.send(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data format is not initialized")
}

func TestBaseGateway_Send_NilProtocolHandler(t *testing.T) {
	bg := newBaseGateway[string, int]("test", &mockSerde{}, nil, backoff.RetryConfig{})

	_, err := bg.send(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "protocol handler is not initialized")
}

func TestBaseGateway_Send_SerializeError(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, backoff.RetryConfig{})

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(errors.New("serialize error"))

	_, err := bg.send(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error marshaling data")
	mockSerde.AssertExpectations(t)
}

func TestBaseGateway_Send_NilResponse(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, backoff.RetryConfig{})

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(nil)
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte(nil), nil)

	_, err := bg.send(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "received nil response")
	mockSerde.AssertExpectations(t)
	mockProto.AssertExpectations(t)
}

func TestBaseGateway_Send_DeserializeError(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, backoff.RetryConfig{})

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(nil)
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte("response"), nil)
	mockSerde.On("Deserialize", mock.Anything, mock.Anything).Return(errors.New("deserialize error"))

	_, err := bg.send(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error unmarshaling response")
	mockSerde.AssertExpectations(t)
	mockProto.AssertExpectations(t)
}

func TestBaseGateway_SendWithRetry_DeadlineExceeded(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, retryConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(nil)
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte{}, context.DeadlineExceeded)

	_, err := bg.sendWithRetry(ctx, "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operation cancelled or timed out")
	mockSerde.AssertExpectations(t)
	mockProto.AssertExpectations(t)
}

func TestBaseGateway_SendWithRetry_MaxAttemptsReached(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, retryConfig)

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(nil)
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte{}, errors.New("some error")).Times(4)

	_, err := bg.sendWithRetry(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max retries reached")
	mockSerde.AssertExpectations(t)
	mockProto.AssertExpectations(t)
	mockProto.AssertNumberOfCalls(t, "Send", 4) // Initial attempt + 3 retries
}

func TestBaseGateway_SendWithRetry_Success(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, retryConfig)

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(nil)
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte("42"), nil)
	mockSerde.On("Deserialize", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		*args.Get(1).(*int) = 42
	}).Return(nil)

	result, err := bg.sendWithRetry(context.Background(), "test_data")

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
	mockSerde.AssertExpectations(t)
	mockProto.AssertExpectations(t)
	mockProto.AssertNumberOfCalls(t, "Send", 1)
}

func TestBaseGateway_SendWithRetry_SuccessAfterRetry(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, retryConfig)

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(nil)
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte{}, errors.New("error")).Twice()
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte("42"), nil).Once()
	mockSerde.On("Deserialize", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		*args.Get(1).(*int) = 42
	}).Return(nil)

	result, err := bg.sendWithRetry(context.Background(), "test_data")

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
	mockSerde.AssertExpectations(t)
	mockProto.AssertExpectations(t)
	mockProto.AssertNumberOfCalls(t, "Send", 3)
}

func TestBaseGateway_SendWithRetry_GatewayUnavailable(t *testing.T) {
	mockSerde := &mockSerde{}
	mockProto := &mockProtocol{}
	retryConfig := backoff.RetryConfig{
		MaxRetries: 3,
		Backoff:    backoff.NewExponentialBackoff(100*time.Millisecond, 1.2, 1*time.Second),
	}
	bg := newBaseGateway[string, int]("test", mockSerde, mockProto, retryConfig)

	mockSerde.On("Serialize", mock.Anything, mock.Anything).Return(nil)
	mockProto.On("Send", mock.Anything, mock.Anything).Return([]byte{}, ErrGatewayUnavailable).Times(4)

	_, err := bg.sendWithRetry(context.Background(), "test_data")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max retries reached")
	mockSerde.AssertExpectations(t)
	mockProto.AssertExpectations(t)
	mockProto.AssertNumberOfCalls(t, "Send", 4) // Initial attempt + 3 retries
}

func TestBaseGateway_Name(t *testing.T) {
	bg := newBaseGateway[string, int]("test_gateway", nil, nil, backoff.RetryConfig{})
	assert.Equal(t, "test_gateway", bg.Name())

	bgNoName := newBaseGateway[string, int]("", nil, nil, backoff.RetryConfig{})
	assert.Equal(t, "unnamed gateway", bgNoName.Name())
}
