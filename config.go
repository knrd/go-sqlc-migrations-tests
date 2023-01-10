package main

import "fmt"

// Postgresql single DB connection.
type Postgresql struct {
	Host     string
	Port     int
	DBName   string
	User     string
	Password string
}

func (p *Postgresql) CreateDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", p.Host, p.Port, p.User, p.DBName, p.Password)
}
