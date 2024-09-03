package protocol

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// HTTPProtocol is a protocol handler for HTTP connections.
type HTTPProtocol struct {
	Client *http.Client
	URL    string
	Method string
}

func NewHTTPConnection(client *http.Client, method, url string) *HTTPProtocol {
	return &HTTPProtocol{
		Client: client,
		Method: method,
		URL:    url,
	}
}

func (h *HTTPProtocol) Send(ctx context.Context, data []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, h.Method, h.URL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to close response body", "error", err)
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	return body, nil
}
