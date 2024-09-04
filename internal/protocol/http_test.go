package protocol

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPProtocol_Send(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		serverResponse string
		serverStatus   int
		requestBody    string
		expectedError  bool
	}{
		{
			name:           "Successful POST request",
			method:         http.MethodPost,
			serverResponse: "response data",
			serverStatus:   http.StatusOK,
			requestBody:    "test data",
			expectedError:  false,
		},
		{
			name:           "Successful GET request",
			method:         http.MethodGet,
			serverResponse: "get response",
			serverStatus:   http.StatusOK,
			requestBody:    "",
			expectedError:  false,
		},
		{
			name:           "Server error",
			method:         http.MethodPost,
			serverResponse: "internal server error",
			serverStatus:   http.StatusInternalServerError,
			requestBody:    "test data",
			expectedError:  true,
		},
		{
			name:           "Not Found error",
			method:         http.MethodGet,
			serverResponse: "not found",
			serverStatus:   http.StatusNotFound,
			requestBody:    "",
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Errorf("Expected %s request, got %s", tt.method, r.Method)
				}

				if tt.requestBody != "" {
					body := make([]byte, r.ContentLength)
					_, _ = r.Body.Read(body)
					if string(body) != tt.requestBody {
						t.Errorf("Expected request body '%s', got '%s'", tt.requestBody, string(body))
					}
				}

				w.WriteHeader(tt.serverStatus)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			httpProtocol := NewHTTPConnection(client, tt.method, server.URL)

			ctx := context.Background()
			response, err := httpProtocol.Send(ctx, []byte(tt.requestBody))

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if string(response) != tt.serverResponse {
					t.Errorf("Expected response '%s', got '%s'", tt.serverResponse, string(response))
				}
			}
		})
	}
}

func TestHTTPProtocol_SendWithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("delayed response"))
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	httpProtocol := NewHTTPConnection(client, http.MethodGet, server.URL)

	ctx := context.Background()
	_, err := httpProtocol.Send(ctx, nil)

	if err == nil {
		t.Errorf("Expected a timeout error, but got none")
	}
}

func TestHTTPProtocol_SendWithInvalidURL(t *testing.T) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	httpProtocol := NewHTTPConnection(client, http.MethodGet, "http://invalid-url")

	ctx := context.Background()
	_, err := httpProtocol.Send(ctx, nil)

	if err == nil {
		t.Errorf("Expected an error for invalid URL, but got none")
	}
}
