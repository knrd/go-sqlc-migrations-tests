package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/knrd/go-sqlc-migrations-tests/database/sqlc_models"
)

// func commit(db *sql.DB) {
// 	if _, err := db.Exec("COMMIT"); err != nil {
// 		log.Fatal("failed to commit: %w", err)
// 	}
// }

type typeAdminConnectionStr string

func (ta typeAdminConnectionStr) String() string {
	return string(ta)
}

type typeTestUserConnectionStr string

func (tu typeTestUserConnectionStr) String() string {
	return string(tu)
}

// set up test schema
func TestingDBSetup(connection typeAdminConnectionStr) error {
	con, err := sql.Open("postgres", string(connection))
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer con.Close()

	if _, err := con.Exec("CREATE DATABASE " + testDbConfig.DBName); err != nil {
		return fmt.Errorf("failed to create DATABASE %v: %w", testDbConfig.DBName, err)
	}
	if _, err := con.Exec(fmt.Sprintf(
		"CREATE USER %s WITH PASSWORD '%s'",
		testDbConfig.User,
		testDbConfig.Password)); err != nil {
		return fmt.Errorf("failed to create USER %v: %w", testDbConfig.User, err)
	}
	if _, err := con.Exec(fmt.Sprintf(
		"GRANT ALL PRIVILEGES ON DATABASE %s TO %s; ALTER DATABASE %s OWNER TO %s;",
		testDbConfig.DBName,
		testDbConfig.User,
		testDbConfig.DBName,
		testDbConfig.User)); err != nil {
		return fmt.Errorf("failed to grant all privileges on database: %w", err)
	}

	var transactionIsolation string
	expectedTransactionIsolationLevel := "read committed"

	row := con.QueryRow("SHOW TRANSACTION ISOLATION LEVEL")
	if err := row.Scan(&transactionIsolation); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
	}

	if transactionIsolation != expectedTransactionIsolationLevel {
		return fmt.Errorf(
			"SHOW TRANSACTION ISOLATION LEVEL expected '%s', got '%s'",
			expectedTransactionIsolationLevel,
			transactionIsolation)
	}
	return nil
}

func getMigrationsDriver(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	path, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path+"/database/migrations",
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to setup migration: %w", err)
	}

	return m, nil
}

// create test tables
func TestingTableCreate(connection typeTestUserConnectionStr) error {
	db, err := sql.Open("postgres", string(connection))
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	migration, err := getMigrationsDriver(db)
	if err != nil {
		return fmt.Errorf("failed to setup migrations driver: %w", err)
	}
	defer migration.Close()

	// or m.Step(2) if you want to explicitly set the number of migrations to run
	if err := migration.Up(); err != nil {
		return fmt.Errorf("failed UP migrations: %w", err)
	}

	return nil
}

func TestingTableNoLeftovers(connection typeTestUserConnectionStr) error {
	db, err := sql.Open("postgres", string(connection))
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	queries := sqlc_models.New(db)
	if num, err := queries.BalancesCount(ctx); num != 0 || err != nil {
		return fmt.Errorf("queries.BalancesCount got %d; err: %w", num, err)
	}
	if num, err := queries.BalanceLogsCount(ctx); num != 0 || err != nil {
		return fmt.Errorf("queries.BalanceLogsCount got %d; err: %w", num, err)
	}
	return nil
}

func TestingDBTeardown(connection typeTestUserConnectionStr) error {
	db, err := sql.Open("postgres", string(connection))
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	migration, err := getMigrationsDriver(db)
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed DOWN migrations: %w", err)
	}
	return nil
}

func TestingDBDropDatabaseAndUser(connection typeAdminConnectionStr) error {
	db, err := sql.Open("postgres", string(connection))
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	if _, err := db.Exec("DROP DATABASE IF EXISTS " + testDbConfig.DBName); err != nil {
		return fmt.Errorf("failed to drop DATABASE: %w", err)
	}
	if _, err := db.Exec("DROP USER IF EXISTS " + testDbConfig.User); err != nil {
		return fmt.Errorf("failed to drop USER: %w", err)
	}
	return nil
}

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var ContextKeyDb = contextKey("db")
var ContextKeyTx = contextKey("tx")

// TestSetupTx create tx and cleanup func for test
func TestSetupTx(t *testing.T, txOptions *sql.TxOptions) (*sqlc_models.Queries, context.Context, func()) {
	t.Helper()

	defaultDb := testDbConfig
	ctx := context.Background()

	db, err := sql.Open("postgres", defaultDb.CreateDSN())
	if err != nil {
		t.Fatal(err)
	}
	queries := sqlc_models.New(db)

	tx, err := db.BeginTx(ctx, txOptions)
	if err != nil {
		t.Fatal(err)
	}
	ctx = context.WithValue(ctx, ContextKeyDb, db)
	ctx = context.WithValue(ctx, ContextKeyTx, tx)
	qtx := queries.WithTx(tx)

	cleanup := func() {
		tx.Rollback()
		db.Close()
	}
	return qtx, ctx, cleanup
}

func TestSetupSubTxFromContext(t *testing.T, ctx context.Context, scTxQueries *sqlc_models.Queries, txOptions *sql.TxOptions) (*sqlc_models.Queries, *sql.Tx, func()) {
	t.Helper()
	db := ctx.Value(ContextKeyDb).(*sql.DB)

	txInner, err := db.BeginTx(ctx, txOptions)
	if err != nil {
		t.Fatal(err)
	}
	qtxInner := scTxQueries.WithTx(txInner)

	cleanup := func() {
		txInner.Rollback()
	}

	return qtxInner, txInner, cleanup
}

// TODO: implement Savepoint using tx from ContextKeyTx
// https://www.postgresql.org/docs/current/sql-savepoint.html
