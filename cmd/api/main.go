package main

import (
	"context"
	"fmt"

	"github.com/rauf/payment-service/internal/gateway"
	"github.com/rauf/payment-service/internal/models"
	"github.com/rauf/payment-service/internal/router"
	"github.com/rauf/payment-service/internal/service"
)

func main() {

	registry := gateway.NewGatewayRegistry()
	registry.Register("gateway-a", gateway.NewGatewayA("http://gateway-a.com"))

	paymentService := service.NewPaymentService(router.NewRouter(registry))
	paymentService.Deposit(context.Background(), models.DepositRequest{
		Amount:   100,
		Currency: "USD",
	})
	fmt.Println("Deposit successful")
}
