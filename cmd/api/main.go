package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rauf/payment-service/cmd/api/handlers"
	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/router"
	"github.com/rauf/payment-service/internal/service"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to run server", "error", err)
	}
}

func run(ctx context.Context) error {
	app, err := setupApplication()
	if err != nil {
		return fmt.Errorf("failed to setup application: %w", err)
	}

	mux, err := app.SetupRoutes()
	if err != nil {
		return fmt.Errorf("failed to setup routes: %w", err)
	}
	slog.InfoContext(ctx, "starting server on :8080")
	return http.ListenAndServe(":8080", mux)
}

func setupApplication() (*Application, error) {
	registry, err := createGatewayRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway registry: %w", err)
	}
	r := router.NewRouter(registry)
	paymentService := service.NewPaymentService(r)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	return NewApplication(registry, paymentHandler), nil
}

func createGatewayRegistry() (*gateway.Registry, error) {
	registry := gateway.NewGatewayRegistry()
	err := registry.Register("Gateway-A", gateway.NewGatewayA("http://gateway-a.com"))
	if err != nil {
		return nil, fmt.Errorf("failed to register Gateway-A: %w", err)
	}
	return registry, nil
}
