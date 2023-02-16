-- name: BalancesSelectAll :many
SELECT * FROM balances;

-- name: BalancesDeleteAll :exec
DELETE FROM balances;

-- name: BalancesCount :one
SELECT COUNT(*) FROM balances;

-- name: BalanceSelectByEmail :one
SELECT * FROM balances WHERE email=$1;

-- name: BalanceSelectForUpdateByEmail :one
SELECT * FROM balances WHERE email=$1 FOR UPDATE;

-- name: BalanceUpdateFByEmail :exec
UPDATE balances SET amount = amount + $1 WHERE email=$2;

-- name: BalanceInsert :exec
INSERT INTO balances (amount, email, created_at)
VALUES ($1, $2, $3);

-- name: BalanceLogsSelectAll :many
SELECT * FROM balance_logs;

-- name: BalanceLogsDeleteAll :exec
DELETE FROM balance_logs;

-- name: BalanceLogsCount :one
SELECT COUNT(*) FROM balance_logs;

-- name: BalanceLogsSelectById :one
SELECT * FROM balance_logs WHERE id=$1;

-- name: BalanceLogsSelectAllByBalanceId :many
SELECT * FROM balance_logs WHERE balance_id=$1;

-- name: BalanceLogsInsert :one
INSERT INTO balance_logs (balance_id, balance_before_change, change, note, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
