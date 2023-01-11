package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// CREATE USER db_tests; -- this user is for test

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

	if _, err := con.Exec("CREATE DATABASE db_tests"); err != nil {
		return fmt.Errorf("failed to create DATABASE db_tests: %w", err)
	}
	if _, err := con.Exec("CREATE USER db_tests_user WITH PASSWORD 'test1234'"); err != nil {
		return fmt.Errorf("failed to create USER db_tests_user: %w", err)
	}
	if _, err := con.Exec("GRANT ALL PRIVILEGES ON DATABASE db_tests TO db_tests_user; ALTER DATABASE db_tests OWNER TO db_tests_user;"); err != nil {
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
	if _, err := db.Exec("DROP DATABASE IF EXISTS db_tests"); err != nil {
		return fmt.Errorf("failed to drop DATABASE: %w", err)
	}
	if _, err := db.Exec("DROP USER IF EXISTS db_tests_user"); err != nil {
		return fmt.Errorf("failed to drop USER: %w", err)
	}
	return nil
}
