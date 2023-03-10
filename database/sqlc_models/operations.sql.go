// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: operations.sql

package sqlc_models

import (
	"context"
	"time"
)

const balanceInsert = `-- name: BalanceInsert :exec
INSERT INTO balances (amount, email, created_at)
VALUES ($1, $2, $3)
`

type BalanceInsertParams struct {
	Amount    int32
	Email     string
	CreatedAt time.Time
}

func (q *Queries) BalanceInsert(ctx context.Context, arg BalanceInsertParams) error {
	_, err := q.db.ExecContext(ctx, balanceInsert, arg.Amount, arg.Email, arg.CreatedAt)
	return err
}

const balanceLogsCount = `-- name: BalanceLogsCount :one
SELECT COUNT(*) FROM balance_logs
`

func (q *Queries) BalanceLogsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, balanceLogsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const balanceLogsDeleteAll = `-- name: BalanceLogsDeleteAll :exec
DELETE FROM balance_logs
`

func (q *Queries) BalanceLogsDeleteAll(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, balanceLogsDeleteAll)
	return err
}

const balanceLogsInsert = `-- name: BalanceLogsInsert :one
INSERT INTO balance_logs (balance_id, balance_before_change, change, note, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type BalanceLogsInsertParams struct {
	BalanceID           int32
	BalanceBeforeChange int32
	Change              int32
	Note                string
	CreatedAt           time.Time
}

func (q *Queries) BalanceLogsInsert(ctx context.Context, arg BalanceLogsInsertParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, balanceLogsInsert,
		arg.BalanceID,
		arg.BalanceBeforeChange,
		arg.Change,
		arg.Note,
		arg.CreatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const balanceLogsSelectAll = `-- name: BalanceLogsSelectAll :many
SELECT id, balance_id, balance_before_change, change, note, created_at FROM balance_logs
`

func (q *Queries) BalanceLogsSelectAll(ctx context.Context) ([]BalanceLog, error) {
	rows, err := q.db.QueryContext(ctx, balanceLogsSelectAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BalanceLog
	for rows.Next() {
		var i BalanceLog
		if err := rows.Scan(
			&i.ID,
			&i.BalanceID,
			&i.BalanceBeforeChange,
			&i.Change,
			&i.Note,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const balanceLogsSelectAllByBalanceId = `-- name: BalanceLogsSelectAllByBalanceId :many
SELECT id, balance_id, balance_before_change, change, note, created_at FROM balance_logs WHERE balance_id=$1
`

func (q *Queries) BalanceLogsSelectAllByBalanceId(ctx context.Context, balanceID int32) ([]BalanceLog, error) {
	rows, err := q.db.QueryContext(ctx, balanceLogsSelectAllByBalanceId, balanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BalanceLog
	for rows.Next() {
		var i BalanceLog
		if err := rows.Scan(
			&i.ID,
			&i.BalanceID,
			&i.BalanceBeforeChange,
			&i.Change,
			&i.Note,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const balanceLogsSelectById = `-- name: BalanceLogsSelectById :one
SELECT id, balance_id, balance_before_change, change, note, created_at FROM balance_logs WHERE id=$1
`

func (q *Queries) BalanceLogsSelectById(ctx context.Context, id int64) (BalanceLog, error) {
	row := q.db.QueryRowContext(ctx, balanceLogsSelectById, id)
	var i BalanceLog
	err := row.Scan(
		&i.ID,
		&i.BalanceID,
		&i.BalanceBeforeChange,
		&i.Change,
		&i.Note,
		&i.CreatedAt,
	)
	return i, err
}

const balanceSelectByEmail = `-- name: BalanceSelectByEmail :one
SELECT id, amount, email, created_at FROM balances WHERE email=$1
`

func (q *Queries) BalanceSelectByEmail(ctx context.Context, email string) (Balance, error) {
	row := q.db.QueryRowContext(ctx, balanceSelectByEmail, email)
	var i Balance
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const balanceSelectForUpdateByEmail = `-- name: BalanceSelectForUpdateByEmail :one
SELECT id, amount, email, created_at FROM balances WHERE email=$1 FOR UPDATE
`

func (q *Queries) BalanceSelectForUpdateByEmail(ctx context.Context, email string) (Balance, error) {
	row := q.db.QueryRowContext(ctx, balanceSelectForUpdateByEmail, email)
	var i Balance
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const balanceUpdateFByEmail = `-- name: BalanceUpdateFByEmail :exec
UPDATE balances SET amount = amount + $1 WHERE email=$2
`

type BalanceUpdateFByEmailParams struct {
	Amount int32
	Email  string
}

func (q *Queries) BalanceUpdateFByEmail(ctx context.Context, arg BalanceUpdateFByEmailParams) error {
	_, err := q.db.ExecContext(ctx, balanceUpdateFByEmail, arg.Amount, arg.Email)
	return err
}

const balancesCount = `-- name: BalancesCount :one
SELECT COUNT(*) FROM balances
`

func (q *Queries) BalancesCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, balancesCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const balancesDeleteAll = `-- name: BalancesDeleteAll :exec
DELETE FROM balances
`

func (q *Queries) BalancesDeleteAll(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, balancesDeleteAll)
	return err
}

const balancesSelectAll = `-- name: BalancesSelectAll :many
SELECT id, amount, email, created_at FROM balances
`

func (q *Queries) BalancesSelectAll(ctx context.Context) ([]Balance, error) {
	rows, err := q.db.QueryContext(ctx, balancesSelectAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Balance
	for rows.Next() {
		var i Balance
		if err := rows.Scan(
			&i.ID,
			&i.Amount,
			&i.Email,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
