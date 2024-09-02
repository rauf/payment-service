-- name: GetAll :one
SELECT *
FROM transaction;

-- name: CreateTransaction :exec
INSERT INTO transaction (type,
                     amount,
                     currency,
                     payment_method,
                     description,
                     customer_id,
                     gateway,
                     gateway_ref_id,
                     status,
                     preferred_gateway,
                     metadata)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11);


-- name: UpdateTransactionStatus :exec
UPDATE transaction
SET status = $1
WHERE gateway_ref_id = $2 AND gateway = $3;

-- name: GetTransactionByGatewayRefId :one
SELECT *
FROM transaction
WHERE gateway_ref_id = $1 AND gateway = $2;