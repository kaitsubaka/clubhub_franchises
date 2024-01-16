package graph

import (
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/db"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/queue"
	ucport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/usecase"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/usecase"
	"github.com/kaitsubaka/clubhub_franchises/internal/infra/repository/psql"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	queue                   *queue.Queue
	franchiseGRCPConn       *grpc.ClientConn
	pendingFranchizeUseCase ucport.PendingFranchiseUseCase
	franchiseUseCase        ucport.FranchiseUseCase
}

func NewResolver() *Resolver {
	db, err := gorm.Open(postgres.Open(db.NewConnectionString(db.PostgresOptions{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     "postgres",
		Password: "admin",
		DBName:   "franchises_db",
		Port:     os.Getenv("POSTGRES_PORT"),
	})), new(gorm.Config))
	if err != nil {
		panic(fmt.Errorf("graph.NewResolver: error creating db connection (%w)", err))
	}
	pendingFranchiseRepository := psql.NewPendingFranchizeRepository(db)
	franchiseGRCPConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", os.Getenv("MS_CLIENT_HOST"), os.Getenv("MS_CLIENT_PORT")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Errorf("graph.NewResolver: error creating ms connection (%w)", err))
	}
	q, err := queue.New(100).
		WithNumOfWorker(10).
		WithSubscriber(usecase.NewProcessFranchiseUseCase(franchiseGRCPConn, pendingFranchiseRepository)).
		Build()
	if err != nil {
		panic(fmt.Errorf("graph.NewResolver: error creating db connection (%w)", err))
	}
	return &Resolver{
		queue:             q,
		franchiseGRCPConn: franchiseGRCPConn,
		pendingFranchizeUseCase: usecase.NewPendingFranchiseUseCase(
			pendingFranchiseRepository,
			q,
		),
		franchiseUseCase: usecase.NewFranchiseUseCase(
			nil,
			psql.NewCountryRepository(db),
			psql.NewCityRepository(db),
			psql.NewCompanyRepository(db),
			psql.NewLocationRepository(db),
			psql.NewFranchiseRepository(db),
			psql.NewDetailedFranchiseRepository(db),
		),
	}
}

func (r *Resolver) ShutDown() {
	r.franchiseGRCPConn.Close()
	r.queue.Close()
}
