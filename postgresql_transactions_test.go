package main

import (
	"database/sql"
	"testing"

	"github.com/knrd/go-sqlc-migrations-tests/database/sqlc_models"
)

func assertNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Error(err)
	}
}

func assertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TODO: remove all .Commit() statements! Change them to savepoints!

// Replicating SELECT FOR UPDATE problems on different isolation levels
// https://asciinema.org/a/q5ovDtaoRYaz9wLviruNhxUws?speed=1.2
// Here on "read committed"
func TestBalanceReadCommited(t *testing.T) {
	sqlcQTx, ctx, cleanup := TestSetupTx(t, nil)
	defer cleanup()

	tx := ctx.Value(ContextKeyTx).(*sql.Tx)

	// num, err := sqlcQTx.BalancesCount(ctx)
	// assertCount(t, num, err, 0)
	_ = assertCount

	email := "test123@example.com"
	data := sqlc_models.BalanceInsertParams{
		Amount: 50,
		Email:  email,
	}

	assertNoError(t,
		sqlcQTx.BalanceInsert(ctx, data))
	tx.Commit()

	qtxInner_A, tx_A, _ := TestSetupSubTxFromContext(t, ctx, sqlcQTx, nil)
	qtxInner_B, tx_B, _ := TestSetupSubTxFromContext(t, ctx, sqlcQTx, nil)
	qtxSummary, txSummary, _ := TestSetupSubTxFromContext(t, ctx, sqlcQTx, nil)

	balance_A, err := qtxInner_A.BalanceSelectByEmail(ctx, email)
	assertNoError(t, err)
	balance_A_change := int32(-40)

	balance_B, err := qtxInner_B.BalanceSelectByEmail(ctx, email)
	assertNoError(t, err)
	balance_B_change := int32(-30)

	assertNoError(t,
		qtxInner_A.BalanceUpdateFByEmail(ctx, sqlc_models.BalanceUpdateFByEmailParams{Email: email, Amount: balance_A_change}))
	log_A_ID, err := qtxInner_A.BalanceLogsInsert(ctx,
		sqlc_models.BalanceLogsInsertParams{
			BalanceID:           balance_A.ID,
			BalanceBeforeChange: balance_A.Amount,
			Change:              int32(balance_A_change)})
	assertNoError(t, err)
	assertNoError(t,
		tx_A.Commit())

	assertNoError(t,
		qtxInner_B.BalanceUpdateFByEmail(ctx, sqlc_models.BalanceUpdateFByEmailParams{Email: email, Amount: balance_B_change}))
	log_B_ID, err := qtxInner_B.BalanceLogsInsert(ctx,
		sqlc_models.BalanceLogsInsertParams{
			BalanceID:           balance_B.ID,
			BalanceBeforeChange: balance_B.Amount,
			Change:              int32(balance_B_change)})
	assertNoError(t, err)
	assertNoError(t,
		tx_B.Commit())

	balance, err := qtxSummary.BalanceSelectByEmail(ctx, email)
	assertNoError(t, err)
	expectedBalance := data.Amount + balance_A_change + balance_B_change
	if balance.Amount != expectedBalance {
		t.Fatalf("Got %d, want %d", balance.Amount, expectedBalance)
	}

	log_A, _ := qtxSummary.BalanceLogsSelectById(ctx, log_A_ID)
	log_B, _ := qtxSummary.BalanceLogsSelectById(ctx, log_B_ID)
	assertEqual(t, log_A.BalanceBeforeChange, data.Amount)
	assertEqual(t, log_B.BalanceBeforeChange, data.Amount)

	// fmt.Println(qtxSummary.BalanceLogsSelectAll(ctx))
	// fmt.Println(qtxSummary.BalancesSelectAll(ctx))

	qtxSummary.BalanceLogsDeleteAll(ctx)
	qtxSummary.BalancesDeleteAll(ctx)
	txSummary.Commit()
}

func TestBalanceSerializableSelectForUpdate(t *testing.T) {
	isolationLevel := &sql.TxOptions{Isolation: sql.LevelSerializable}
	sqlcQTx, ctx, cleanup := TestSetupTx(t, isolationLevel)
	defer cleanup()

	tx := ctx.Value(ContextKeyTx).(*sql.Tx)

	// num, err := sqlcQTx.BalancesCount(ctx)
	// assertCount(t, num, err, 0)
	_ = assertCount

	email := "test123@example.com"
	data := sqlc_models.BalanceInsertParams{
		Amount: 50,
		Email:  email,
	}

	assertNoError(t,
		sqlcQTx.BalanceInsert(ctx, data))
	tx.Commit()

	qtxInner_A, tx_A, cleanup_A := TestSetupSubTxFromContext(t, ctx, sqlcQTx, isolationLevel)
	qtxInner_B, _, cleanup_B := TestSetupSubTxFromContext(t, ctx, sqlcQTx, isolationLevel)
	qtxSummary, txSummary, cleanup_Summary := TestSetupSubTxFromContext(t, ctx, sqlcQTx, nil)
	defer cleanup_A()
	defer cleanup_B()
	defer cleanup_Summary()

	_, err := qtxInner_B.BalanceSelectByEmail(ctx, email)
	assertNoError(t, err)
	balance_B_change := int32(-30)

	balance_A, err := qtxInner_A.BalanceSelectByEmail(ctx, email)
	assertNoError(t, err)
	balance_A_change := int32(-40)

	assertNoError(t,
		qtxInner_A.BalanceUpdateFByEmail(ctx, sqlc_models.BalanceUpdateFByEmailParams{Email: email, Amount: balance_A_change}))
	_, err = qtxInner_A.BalanceLogsInsert(ctx,
		sqlc_models.BalanceLogsInsertParams{
			BalanceID:           balance_A.ID,
			BalanceBeforeChange: balance_A.Amount,
			Change:              int32(balance_A_change)})
	assertNoError(t, err)
	assertNoError(t,
		tx_A.Commit())

	err = qtxInner_B.BalanceUpdateFByEmail(ctx, sqlc_models.BalanceUpdateFByEmailParams{Email: email, Amount: balance_B_change})
	assertEqual(t, err.Error(), "pq: could not serialize access due to concurrent update")

	qtxSummary.BalanceLogsDeleteAll(ctx)
	qtxSummary.BalancesDeleteAll(ctx)
	txSummary.Commit()
}
