-- name: BalancesGetAll :many
SELECT * FROM balances;

-- name: BalanceGetByEmail :one
SELECT * FROM balances WHERE email=$1;

-- name: BalanceInsert :exec
INSERT INTO balances (amount, email, created_at)
VALUES ($1, $2, $3);

-- name: BalanceLogInsert :exec
INSERT INTO balance_logs (balance_id, change, note, created_at)
VALUES ($1, $2, $3, $4);
