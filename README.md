## Payment Service

A flexible and extensible payment service designed to integrate with multiple payment gateways. 
You can find the design document [here](design.md) and openapi spec [here](openapi.yaml)

### Features

- Modular architecture for easy integration of new payment gateways
- Configurable retry mechanism with customizable backoff strategies
- Flexible serialization/deserialization support
- Abstracted protocol handling for various communication methods
- Comprehensive error handling and logging

### Prerequisites

- Go 1.23
- make
- Docker

### Build and Run

1. Install necessary dependencies
```bash
make init 
```

2. Run app and postgres
```bash
make docker-compose/up
```

3. Run database migrations
```bash
export PGDATABASE=payment
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/payment
export GOOSE_MIGRATION_DIR=db/migrations
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="$DATABASE_URL"

goose up
```

4. Access the API on http://localhost:8080

To stop API and postgres
```bash
make docker-compose/down
```


### Run Tests

```bash
make test
```


### Curl commands

1. Create transaction

```bash
curl --request POST \
  --url http://localhost:8080/api/v1/transactions \
  --header 'Content-Type: application/json' \
  --data '{
  "amount": 123,
  "type": "withdrawal",
  "currency": "USD",
  "payment_method": "BANK_TRANSFER",
  "description": "payment",
  "customer_id": "cus123",
  "preferred_gateway": "gatewayB",
  "metadata": {
		"orderID": 123
	}
}'
```

2. Gateway A callback

```bash
curl --request POST \
  --url http://localhost:8080/api/v1/gateways/gatewayA/callback \
  --header 'Content-Type: application/json' \
  --data '{
	"ref_id": "TiOIptTKAggASOT5wu3i",
	"status": "success"
}'
```
3. Gateway B callback

```bash
curl --request POST \
  --url http://localhost:8080/api/v1/gateways/gatewayB/callback \
  --header 'Content-Type: application/xml' \
  --data '<?xml version="1.0" encoding="UTF-8"?>
<callback>
  <ref_id>isbKQ4O8LHtIViY61ADa</ref_id>
  <status>success</status>
</callback>'
```

4. Update status API

Replace with id in the path
```bash
curl --request PATCH \
  --url http://localhost:8080/api/v1/transactions/vxOAi2w6ZQB1pilXYitU/status \
  --header 'Content-Type: application/json' \
  --data '{
	"gateway": "gatewayB",
	"status": "success"
}'
```

### Libraries/ Tools Used
1. [sqlc](https://github.com/sqlc-dev/sqlc)
2. [goose](https://github.com/pressly/goose)
3. [sony/gobreaker](https://github.com/sony/gobreaker) circuit breaker

### Further improvements

1. Improve test coverage
2. Take config from environment variables