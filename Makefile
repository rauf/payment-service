
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

.PHONY: docker-compose/up
docker-compose/up:
	@docker compose up -d

.PHONY: docker-compose/down
docker-compose/down:
	@docker compose down

.PHONY: db/migrate/up
# For first time: INSERT INTO goose_db_version (version_id, is_applied) VALUES (0, TRUE);
db/migrate/up:
	@go run ./cmd/migrate up

.PHONY: db/migrate/down
db/migrate/down:
	@go run ./cmd/migrate down

.PHONY: db/migrate/create
# make db/migrate/create name=test123
db/migrate/create:
	@go run ./cmd/migrate create $(name) sql
