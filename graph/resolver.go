package graph

import (
	"fmt"

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
}

func NewResolver() *Resolver {
	db, err := gorm.Open(postgres.Open(db.NewConnectionString(db.PostgresOptions{
		Host:     "localhost",
		User:     "admin",
		Password: "admin",
		DBName:   "test",
		Port:     "5432",
	})), new(gorm.Config))
	if err != nil {
		panic(fmt.Errorf("graph.NewResolver: error creating db connection (%w)", err))
	}
	franchiseGRCPConn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("graph.NewResolver: error creating db connection (%w)", err))
	}
	q, err := queue.New(100).
		WithNumOfWorker(10).
		WithSubscriber(usecase.NewProcessFranchiseUseCase(franchiseGRCPConn)).
		Build()
	if err != nil {
		panic(fmt.Errorf("graph.NewResolver: error creating db connection (%w)", err))
	}
	return &Resolver{
		queue:             q,
		franchiseGRCPConn: franchiseGRCPConn,
		pendingFranchizeUseCase: usecase.NewPendingFranchiseUseCase(
			psql.NewPendingFranchizeRepository(db),
			q,
		),
	}
}

func (r *Resolver) ShutDown() {
	r.franchiseGRCPConn.Close()
	r.queue.Close()
}
