package main

import (
	"log/slog"
	"net/http"

	"github.com/rauf/payment-service/cmd/api/handlers"
	"github.com/rauf/payment-service/internal/utils/jsonutil"
)

func (a *Application) SetupRoutes() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/transact", makeHandler(a.PaymentHandler.HandleTransaction))
	mux.HandleFunc("POST /api/v1/transaction/status", makeHandler(a.PaymentHandler.HandleUpdateStatus))

	return mux, nil
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request) handlers.Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := fn(w, r)
		if res.Err != nil {
			slog.ErrorContext(r.Context(), "error while processing request", "error", res.Err)
			handleError(w, r, res)
			return
		}
		if err := jsonutil.WriteJSON(w, http.StatusOK, res); err != nil {
			slog.Error("failed to encode response", "error", err)
		}
	}
}
func handleError(w http.ResponseWriter, r *http.Request, res handlers.Response) {
	if res.Err != nil {
		if err := jsonutil.WriteJSON(w, res.Code, res); err != nil {
			slog.ErrorContext(r.Context(), "error while processing request", "error", err, "internal_error", res.Err)
			slog.Error("failed to encode error response", "error", err)
		}
		return
	}
	internalErrResponse := &handlers.Response{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	}

	if err := jsonutil.WriteJSON(w, http.StatusInternalServerError, internalErrResponse); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}
