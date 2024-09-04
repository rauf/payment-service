package main

import (
	"net/http"

	"github.com/rauf/payment-service/cmd/api/handlers"
)

func (a *Application) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/transactions", handlers.MakeHandler(a.PaymentHandler.HandleCreateTransaction))
	mux.HandleFunc("PATCH /api/v1/transactions/{id}/status", handlers.MakeHandler(a.PaymentHandler.HandleUpdateStatus))

	// Each gateway can have its own response and format
	mux.HandleFunc("POST /api/v1/gateways/gatewayA/callback", handlers.MakeHandler(a.PaymentHandler.HandleGatewayACallback))
	mux.HandleFunc("POST /api/v1/gateways/gatewayB/callback", handlers.MakeHandler(a.PaymentHandler.HandleGatewayBCallback))

	return mux
}
