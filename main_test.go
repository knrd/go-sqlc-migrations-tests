package main

import (
	"os"
	"testing"

	"log"
)

// TestMain model package setup/teardonw
func TestMain(m *testing.M) {
	// for create/drop schema
	createSchemaCon := adminDbConfig.CreateDSN()
	// for create database objects
	createTableCon := testDbConfig.CreateDSN()

	log.Println("0. TestingDBTeardown(createSchemaCon)")
	if err := TestingDBTeardown(createSchemaCon); err != nil {
		log.Fatal(err)
	}
	log.Println("1. TestingDBSetup(createSchemaCon)")
	if err := TestingDBSetup(createSchemaCon); err != nil {
		log.Fatal(err)
	}
	log.Println("2. TestingTableCreate(createTableCon)")
	if err := TestingTableCreate(createTableCon); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	log.Println("3. TestingDBTeardown(createSchemaCon)")
	if err := TestingDBTeardown(createSchemaCon); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
