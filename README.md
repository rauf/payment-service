## Payment Service

A flexible and extensible payment service designed to integrate with multiple payment gateways. You can find the design document [here](design.md)

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

```bash
make init 
make docker-compose/up
export PGDATABASE=payment
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/payment
export GOOSE_MIGRATION_DIR=db/migrations
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="$DATABASE_URL"

goose up
```

### Libraries/ Tools Used
1. [sqlc](https://github.com/sqlc-dev/sqlc)
2. [goose](https://github.com/pressly/goose)
3. [sony/gobreaker](https://github.com/sony/gobreaker) circuit breaker

### Improvements

1. Improve test coverage
2. Take config from environment variables