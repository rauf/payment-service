package protocol

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rauf/payment-service/internal/utils/randutil"
)

// HTTPProtocolMock is a mock implementation of the HTTP protocol. To be used for testing purposes
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
		RefID     string    `json:"ref_id"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}{
		RefID:     randutil.RandomString(20),
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	return json.Marshal(res)
}
