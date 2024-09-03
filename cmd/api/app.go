package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rauf/payment-service/cmd/api/handlers"
	"github.com/rauf/payment-service/internal/backoff"
	"github.com/rauf/payment-service/internal/config"
	"github.com/rauf/payment-service/internal/consts"
	"github.com/rauf/payment-service/internal/database"
	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/registry"
	"github.com/rauf/payment-service/internal/repo"
	"github.com/rauf/payment-service/internal/router"
	"github.com/rauf/payment-service/internal/service"
	"github.com/sony/gobreaker/v2"
)

// Application is the main application struct that holds the dependencies.
type Application struct {
	Registry       *registry.Registry[gateway.PaymentGateway]
	PaymentHandler *handlers.PaymentHandler
}

func NewApplication(regis *registry.Registry[gateway.PaymentGateway], ph *handlers.PaymentHandler) *Application {
	return &Application{
		Registry:       regis,
		PaymentHandler: ph,
	}
}

// setupApplication creates a new application instance with the required dependencies.
func setupApplication() (*Application, error) {
	settings := gobreaker.Settings{ // TODO: take input from config
		MaxRequests: 1,
		Interval:    5 * time.Minute,
		Timeout:     5 * time.Minute,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 1
		},
	}

	gatewayRegistry, err := createGatewayRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway registry: %w", err)
	}
	conf := config.NewConfig()
	db, err := database.NewDatabase(conf.Database)
	r := router.NewRouter(gatewayRegistry, settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	paymentRepo := repo.NewPaymentRepo(models.New(db))
	paymentService := service.NewPaymentService(r, paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	return NewApplication(gatewayRegistry, paymentHandler), nil
}

func createGatewayRegistry() (*registry.Registry[gateway.PaymentGateway], error) {
	gatewayRegistry := registry.NewRegistry[gateway.PaymentGateway]()
	httpClient := &http.Client{
		Timeout: 10 * time.Second, // specify the timeout for the http client, should be configurable
	}

	err := gatewayRegistry.Register(consts.GatewayA, gateway.NewGatewayA( // TODO: these should be configurable and read from env variable
		consts.GatewayA,
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
	err = gatewayRegistry.Register(consts.GatewayB, gateway.NewGatewayB(
		consts.GatewayB,
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
	return gatewayRegistry, nil
}
