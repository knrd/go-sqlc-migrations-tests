package main

import (
	"testing"
	"time"

	"github.com/knrd/go-sqlc-migrations-tests/database/sqlc_models"
)

func assertError(t testing.TB, got error, want string) {
	t.Helper()

	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got.Error() != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func assertEqualBalanses(t testing.TB, got []sqlc_models.Balance, want []sqlc_models.BalanceInsertParams, err error) {
	t.Helper()

	if err != nil {
		t.Fatal("Error", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for k, v := range got {
		if v.Amount != want[k].Amount {
			t.Errorf("Amount got %v, want %v", got[k], want[k])
		}
		if v.Email != want[k].Email {
			t.Errorf("Email got %v, want %v", got[k], want[k])
		}
		// t.Log(v.CreatedAt, v.CreatedAt.Format(time.RFC3339), v.CreatedAt.Location())
		// t.Log(want[k].CreatedAt, want[k].CreatedAt.Format(time.RFC3339), want[k].CreatedAt.Location())
		// if v.CreatedAt.Format(time.RFC3339) != want[k].CreatedAt.Format(time.RFC3339) {
		// 	t.Errorf("CreatedAt got %s, want %s", got[k].CreatedAt.Format(time.RFC3339), want[k].CreatedAt.Format(time.RFC3339))
		// }
	}
}

func assertCount(t testing.TB, got int64, err error, want int64) {
	t.Helper()

	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("Count got %v, want %v", got, want)
	}
}

func TestEmailIsRequiredInParallelSubTx(t *testing.T) {
	sqlcQTx, ctx, cleanup := TestSetupTx(t, nil)
	t.Cleanup(func() {
		cleanup()
	})

	wrong_emails := []struct {
		text string
	}{
		{"aaa"},
		{"+test@example.com"},
		{"example.com"},
		{"@"},
	}

	for _, wrong_email := range wrong_emails {
		wrong_email := wrong_email
		t.Run(wrong_email.text, func(t *testing.T) {
			t.Parallel()
			qtxInner, _, cleanupInner := TestSetupSubTxFromContext(t, ctx, sqlcQTx, nil)

			assertError(
				t,
				qtxInner.BalanceInsert(ctx, sqlc_models.BalanceInsertParams{Email: wrong_email.text}),
				"pq: new row for relation \"balances\" violates check constraint \"proper_email\"",
			)

			cleanupInner()
		})
	}
}

func TestBalanceInserted(t *testing.T) {
	sqlcQTx, ctx, cleanup := TestSetupTx(t, nil)
	defer cleanup()

	now := time.Now().UTC().Round(time.Microsecond)
	expected_data := []sqlc_models.BalanceInsertParams{{
		Amount:    11,
		Email:     "test10@example.com",
		CreatedAt: now,
	}}

	if err := sqlcQTx.BalanceInsert(ctx, expected_data[0]); err != nil {
		t.Fatal(err)
	}

	balances, err := sqlcQTx.BalancesSelectAll(ctx)
	assertEqualBalanses(t, balances, expected_data, err)
}

func TestBalanceInsertedWithSubTx(t *testing.T) {
	sqlcQTx, ctx, cleanup := TestSetupTx(t, nil)
	defer cleanup()

	qtxInner1, _, cleanupInner1 := TestSetupSubTxFromContext(t, ctx, sqlcQTx, nil)
	qtxInner2, _, cleanupInner2 := TestSetupSubTxFromContext(t, ctx, sqlcQTx, nil)
	defer cleanupInner1()
	defer cleanupInner2()

	expected_data1 := []sqlc_models.BalanceInsertParams{{
		Amount:    123,
		Email:     "test123@example.com",
		CreatedAt: time.Now().UTC().Round(time.Microsecond),
	}}

	if err := qtxInner1.BalanceInsert(ctx, expected_data1[0]); err != nil {
		t.Fatal(err)
	}
	balances1, err := qtxInner1.BalancesSelectAll(ctx)
	assertEqualBalanses(t, balances1, expected_data1, err)

	expected_data2 := []sqlc_models.BalanceInsertParams{{
		Amount:    22,
		Email:     "test22@example.com",
		CreatedAt: time.Now().UTC().Round(time.Microsecond),
	}}

	if err := qtxInner2.BalanceInsert(ctx, expected_data2[0]); err != nil {
		t.Fatal(err)
	}
	balances2, err := qtxInner2.BalancesSelectAll(ctx)
	assertEqualBalanses(t, balances2, expected_data2, err)

	balances, err := sqlcQTx.BalancesSelectAll(ctx)
	assertEqualBalanses(t, balances, []sqlc_models.BalanceInsertParams{}, err)
}

// Replicating SELECT FOR UPDATE problems on different isolation levels
// https://asciinema.org/a/q5ovDtaoRYaz9wLviruNhxUws?speed=1.2
// Here on "read committed"
