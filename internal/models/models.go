package models

type DepositRequest struct {
	Amount           float64
	Currency         string
	PreferredGateway string
}

type DepositResponse struct {
	TransactionID string
}

type WithdrawalRequest struct {
	Amount           float64
	Currency         string
	PreferredGateway string
}

type WithdrawalResponse struct {
	TransactionID string
}
