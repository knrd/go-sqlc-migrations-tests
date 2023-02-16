package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/knrd/go-sqlc-migrations-tests/database/sqlc_models"

	_ "github.com/lib/pq"
)

func run() error {
	defaultDb := adminDbConfig

	fmt.Println(defaultDb)

	ctx := context.Background()

	fmt.Println(defaultDb.CreateDSN())

	db, err := sql.Open("postgres", defaultDb.CreateDSN())
	if err != nil {
		return err
	}

	queries := sqlc_models.New(db)
	_ = queries

	// list all authors
	authors, err := queries.BalancesSelectAll(ctx)
	if err != nil {
		return err
	}
	for k, v := range authors {
		log.Printf("%d. %+v", k, v)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
