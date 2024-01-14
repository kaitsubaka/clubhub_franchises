package db

import "fmt"

type PostgresOptions struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

func NewConnectionString(opts PostgresOptions) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=america/bogota",
		"localhost",
		"admin",
		"admin",
		"test",
		"5432",
	)
}
