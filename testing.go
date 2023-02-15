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
	_ "github.com/lib/pq"

	"github.com/knrd/go-sqlc-migrations-tests/database/sqlc_models"
)

// func commit(db *sql.DB) {
// 	if _, err := db.Exec("COMMIT"); err != nil {
// 		log.Fatal("failed to commit: %w", err)
// 	}
// }

// TestingDBSetup set up test schema
func TestingDBSetup(conStr string) error {
	con, err := sql.Open("postgres", conStr)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer con.Close()

	if _, err := con.Exec("CREATE DATABASE " + testDbConfig.DBName); err != nil {
		return fmt.Errorf("failed to create DATABASE %v: %w", testDbConfig.DBName, err)
	}
	if _, err := con.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", testDbConfig.User, testDbConfig.Password)); err != nil {
		return fmt.Errorf("failed to create USER %v: %w", testDbConfig.User, err)
	}
	if _, err := con.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s; ALTER DATABASE %s OWNER TO %s;", testDbConfig.DBName, testDbConfig.User, testDbConfig.DBName, testDbConfig.User)); err != nil {
		return fmt.Errorf("failed to grant all privileges on database: %w", err)
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

// TestingTableCreate create test tables
func TestingTableCreate(conStr string) error {
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	m, err := getMigrationsDriver(db)
	if err != nil {
		return fmt.Errorf("failed to setup migrations driver: %w", err)
	}
	defer m.Close()

	// or m.Step(2) if you want to explicitly set the number of migrations to run
	if err := m.Up(); err != nil {
		return fmt.Errorf("failed UP migrations: %w", err)
	}

	return nil
}

// TestingDBTeardown drop test schema
func TestingDBTeardown(conStr string) error {
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	m, err := getMigrationsDriver(db)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed DOWN migrations: %w", err)
	}

	// we can use here migrate.Down() if necessary
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

// var ContextKeyTx = contextKey("tx")

// TestSetupTx create tx and cleanup func for test
func TestSetupTx(t *testing.T) (*sqlc_models.Queries, context.Context, func()) {
	t.Helper()

	defaultDb := testDbConfig
	ctx := context.Background()

	db, err := sql.Open("postgres", defaultDb.CreateDSN())
	if err != nil {
		t.Fatal(err)
	}
	queries := sqlc_models.New(db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	ctx = context.WithValue(ctx, ContextKeyDb, db)
	qtx := queries.WithTx(tx)

	cleanup := func() {
		tx.Rollback()
		db.Close()
	}
	return qtx, ctx, cleanup
}

func TestSetupSubTxFromContext(t *testing.T, ctx context.Context, qtx *sqlc_models.Queries) (*sqlc_models.Queries, *sql.Tx, func()) {
	t.Helper()
	db := ctx.Value(ContextKeyDb).(*sql.DB)
	txInner, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtxInner := qtx.WithTx(txInner)

	cleanup := func() {
		txInner.Rollback()
	}

	return qtxInner, txInner, cleanup
}
