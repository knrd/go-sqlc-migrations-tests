// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package sqlc_models

import (
	"time"
)

type Balance struct {
	ID        int32
	Amount    int32
	Email     string
	CreatedAt time.Time
}

type BalanceLog struct {
	ID                  int64
	BalanceID           int32
	BalanceBeforeChange int32
	Change              int32
	Note                string
	CreatedAt           time.Time
}
