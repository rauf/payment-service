
Steps

```bash
make docker-compose/up

```


## Database migrations

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

## Libraries/ Tools Used
1. sqlc
2. goose