package test

import (
	"context"
	"testing"
	"time"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	pggorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPSQLTestDB(initScript string, t *testing.T) *gorm.DB {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		postgres.WithInitScripts(initScript),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("admin"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}
	// Get the container's host and port
	mariadbHost, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	//obtaining the externally mapped port for the container
	mariadbPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	gormdb, err := gorm.Open(pggorm.Open(db.NewConnectionString(db.PostgresOptions{
		Host:     mariadbHost,
		User:     "postgres",
		Password: "admin",
		DBName:   "franchises_db",
		Port:     mariadbPort.Port(),
	})), new(gorm.Config))
	if err != nil {
		t.Fatal(err)
	}
	return gormdb
}
