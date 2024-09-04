-- +goose Up
-- +goose StatementBegin
CREATE TYPE transaction_status AS ENUM ('PENDING', 'SUCCESS', 'FAILED');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE transaction_type AS ENUM ('DEPOSIT', 'WITHDRAWAL', 'TRANSFER');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction
(
    id                SERIAL PRIMARY KEY,
    type              transaction_type    NOT NULL,
    amount            NUMERIC(15, 2) NOT NULL,
    currency          VARCHAR(10)    NOT NULL,
    payment_method    VARCHAR(50)    NOT NULL,
    description       TEXT,
    customer_id       VARCHAR(100)   NOT NULL,
    gateway           VARCHAR(50)    NOT NULL,
    gateway_ref_id    VARCHAR(50)    NOT NULL,
    status            transaction_status NOT NULL DEFAULT 'PENDING',
    preferred_gateway VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (gateway_ref_id, gateway)
);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS transaction;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS transaction_status;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TYPE IF EXISTS transaction_type;
-- +goose StatementEnd
