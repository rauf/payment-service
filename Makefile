
.PHONY: init
init:
	@go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: gen
gen:
	@sqlc generate

.PHONY: build
build:
	@go build -o ./tmp/api ./cmd/api

.PHONY: run
run:
	@go run ./cmd/api

.PHONY: test
test:
	@go test -v ./...

.PHONY: docker-compose/up
docker-compose/up:
	@docker compose up -d

.PHONY: docker-compose/down
docker-compose/down:
	@docker compose down
