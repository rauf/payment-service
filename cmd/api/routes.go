package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/rauf/payment-service/cmd/api/handlers"
	"github.com/rauf/payment-service/internal/utils/jsonutil"
)

func (a *Application) SetupRoutes() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/deposit", makeHandler(a.PaymentHandler.HandleDeposit))
	mux.HandleFunc("POST /api/v1/withdrawal", makeHandler(a.PaymentHandler.HandleWithdrawal))
	mux.HandleFunc("POST /api/v1/status", makeHandler(a.PaymentHandler.HandleUpdateStatus))

	return mux, nil
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			slog.ErrorContext(r.Context(), "error while processing request", "error", err)
			handleError(w, r, err)
			return
		}
		var apiRes *handlers.Response
		if errors.As(err, &apiRes) {
			if err := jsonutil.WriteJSON(w, http.StatusOK, apiRes.Data); err != nil {
				slog.Error("failed to encode response", "error", err)
			}
		}
	}
}
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	var apiRes *handlers.Response
	if errors.As(err, &apiRes) {
		slog.ErrorContext(r.Context(), "error while processing request", "error", err, "internal_error", apiRes.Err)
		if err := jsonutil.WriteJSON(w, apiRes.Code, apiRes); err != nil {
			slog.Error("failed to encode error response", "error", err)
		}
		return
	}
	res := &handlers.Response{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	}

	if err := jsonutil.WriteJSON(w, http.StatusInternalServerError, res); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}
