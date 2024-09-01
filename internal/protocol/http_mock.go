package protocol

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rauf/payment-service/internal/utils/randutil"
)

type HTTPProtocolMock struct {
	Client *http.Client
	URL    string
	Method string
}

func NewHTTPConnectionMock(client *http.Client, method, url string) *HTTPProtocolMock {
	return &HTTPProtocolMock{
		Client: client,
		Method: method,
		URL:    url,
	}
}

func (h *HTTPProtocolMock) Send(_ context.Context, _ []byte) ([]byte, error) {
	res := struct {
		TransactionID string    `json:"transaction_id"`
		Status        string    `json:"status"`
		CreatedAt     time.Time `json:"created_at"`
	}{
		TransactionID: randutil.RandomString(10),
		Status:        "pending",
		CreatedAt:     time.Now(),
	}
	return json.Marshal(res)
}
