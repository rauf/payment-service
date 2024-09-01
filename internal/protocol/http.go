package protocol

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
)

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
		return nil, err
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to close response body", err)
			return
		}
	}(resp.Body)

	return io.ReadAll(resp.Body)
}
