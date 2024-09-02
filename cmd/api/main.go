package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rauf/payment-service/cmd/api/handlers"
	"github.com/rauf/payment-service/internal/backoff"
	"github.com/rauf/payment-service/internal/config"
	"github.com/rauf/payment-service/internal/database"
	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/repo"
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
	conf := config.NewConfig()
	db, err := database.NewDatabase(conf.Database)
	r := router.NewRouter(registry)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	paymentRepo := repo.NewPaymentRepo(models.New(db))
	paymentService := service.NewPaymentService(r, paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	return NewApplication(registry, paymentHandler), nil
}

func createGatewayRegistry() (*gateway.Registry, error) {
	registry := gateway.NewGatewayRegistry()
	httpClient := &http.Client{
		Timeout: 10 * time.Second, // specify the timeout for the http client, should be configurable
	}

	err := registry.Register("Gateway-A", gateway.NewGatewayA( // TODO: these should be configurable and read from env variable
		"Gateway-A",
		http.MethodPost,
		"http: //gateway-a.com",
		httpClient,
		backoff.RetryConfig{
			MaxRetries: 3,
			Backoff:    backoff.NewExponentialBackoff(1*time.Second, 1.2, 2*time.Second),
		}))
	if err != nil {
		return nil, fmt.Errorf("failed to register Gateway-A: %w", err)
	}
	err = registry.Register("Gateway-B", gateway.NewGatewayB(
		"Gateway-B",
		http.MethodPost,
		"http://gateway-b.com",
		httpClient,
		backoff.RetryConfig{
			MaxRetries: 3,
			Backoff:    backoff.NewExponentialBackoff(1*time.Second, 1.2, 2*time.Second),
		}))
	if err != nil {
		return nil, fmt.Errorf("failed to register Gateway-B: %w", err)
	}
	return registry, nil
}
