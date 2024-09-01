package main

import (
	"github.com/rauf/payment-service/cmd/api/handlers"
	"github.com/rauf/payment-service/internal/gateway"
)

type Application struct {
	Registry       *gateway.Registry
	PaymentHandler *handlers.PaymentHandler
}

func NewApplication(regis *gateway.Registry, ph *handlers.PaymentHandler) *Application {
	return &Application{
		Registry:       regis,
		PaymentHandler: ph,
	}
}
