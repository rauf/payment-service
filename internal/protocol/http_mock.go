package protocol

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/rauf/payment-service/internal/utils/randutil"
)

// HTTPProtocolMock is a mock implementation of the HTTP protocol. To be used for testing purposes
type HTTPProtocolMock struct {
	Client         *http.Client
	URL            string
	Method         string
	ResponseFormat string // "json" or "xml", only for mocking purpose
}

func NewHTTPConnectionMock(client *http.Client, method, url, responseFormat string) *HTTPProtocolMock {
	return &HTTPProtocolMock{
		Client:         client,
		Method:         method,
		URL:            url,
		ResponseFormat: responseFormat,
	}
}

func (h *HTTPProtocolMock) Send(_ context.Context, _ []byte) ([]byte, error) {
	type response struct {
		RefID     string    `json:"ref_id" xml:"ref_id"`
		Status    string    `json:"status" xml:"status"`
		CreatedAt time.Time `json:"created_at" xml:"created_at"`
	}

	res := response{
		RefID:     randutil.RandomString(20),
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	var output []byte
	var err error

	switch h.ResponseFormat {
	case "json":
		output, err = json.Marshal(res)
	case "xml":
		output, err = xml.Marshal(res)
	default:
		return nil, fmt.Errorf("unsupported response format: %s", h.ResponseFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("error marshaling response: %w", err)
	}
	return output, nil
}
