package main

import (
	"os"
	"testing"

	"log"
)

// TestMain model package setup/teardonw
func TestMain(m *testing.M) {
	// for create/drop schema
	administratorConnection := typeAdminConnectionStr(adminDbConfig.CreateDSN())
	// for create database objects
	testUserConnection := typeTestUserConnectionStr(testDbConfig.CreateDSN())

	log.Println("1. TestingDBDropDatabaseAndUser(administratorConnectionStr)")
	if err := TestingDBDropDatabaseAndUser(administratorConnection); err != nil {
		log.Fatal(err)
	}
	log.Println("2. TestingDBSetup(administratorConnectionStr)")
	if err := TestingDBSetup(administratorConnection); err != nil {
		log.Fatal(err)
	}
	log.Println("3. TestingTableCreate(createTableCon)")
	if err := TestingTableCreate(testUserConnection); err != nil {
		log.Fatal(err)
	}

	// now we will run all the tests in the package
	code := m.Run()

	log.Println("4. TestingTableNoLeftovers(createTableCon)")
	if err := TestingTableNoLeftovers(testUserConnection); err != nil {
		log.Fatal(err)
	}
	log.Println("5. TestingDBTeardown(testUserConnectionStr)")
	if err := TestingDBTeardown(testUserConnection); err != nil {
		log.Fatal(err)
	}
	log.Println("6. TestingDBDropDatabaseAndUser(administratorConnectionStr)")
	if err := TestingDBDropDatabaseAndUser(administratorConnection); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
