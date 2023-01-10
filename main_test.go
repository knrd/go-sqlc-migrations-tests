package main

import (
	"flag"
	"os"
	"testing"

	"log"
)

// TestMain model package setup/teardonw
func TestMain(m *testing.M) {
	flag.Parse()

	adminDbConfig := Postgresql{
		Host:     "localhost",
		Port:     15432,
		DBName:   "postgres",
		User:     "postgres",
		Password: "qwerty",
	}

	testDbConfig := Postgresql{
		Host:     "localhost",
		Port:     15432,
		DBName:   "db_tests",
		User:     "db_tests_user",
		Password: "test1234",
	}

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
