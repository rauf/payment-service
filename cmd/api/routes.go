package main

import (
	"log/slog"
	"net/http"
)

func (a *Application) SetupRoutes() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/deposit", makeHandler(a.PaymentHandler.HandleDeposit))

	return mux, nil
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			slog.ErrorContext(r.Context(), "error while processing request", "error", err)
		}
	}
}
