package main

import (
	"net/http"

	"github.com/rauf/payment-service/cmd/api/handlers"
)

func (a *Application) SetupRoutes() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/transact", handlers.MakeHandler(a.PaymentHandler.HandleTransaction))
	mux.HandleFunc("POST /api/v1/transaction/status", handlers.MakeHandler(a.PaymentHandler.HandleUpdateStatus))

	return mux, nil
}
