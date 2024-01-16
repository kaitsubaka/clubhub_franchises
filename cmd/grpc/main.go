package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/common"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/db"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	ucport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/usecase"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/usecase"
	"github.com/kaitsubaka/clubhub_franchises/internal/infra/repository/http"
	"github.com/kaitsubaka/clubhub_franchises/internal/infra/repository/psql"
	pb "github.com/kaitsubaka/clubhub_franchises/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const scrapToolURL = "https://api.ssllabs.com/api/v3/analyze?host=%s"

type FranchiseGRPCServer struct {
	pb.UnimplementedFranchiseServiceServer
	frachiseUseCase ucport.FranchiseUseCase
}

func NewFranchiseGRPCServer(fuc ucport.FranchiseUseCase) *FranchiseGRPCServer {
	return &FranchiseGRPCServer{
		frachiseUseCase: fuc,
	}
}

func (f *FranchiseGRPCServer) Create(ctx context.Context, req *pb.CreateFranchiseRequest) (*pb.SuccessResponse, error) {
	err := f.frachiseUseCase.Create(dto.PendingFranchiseDTO{
		ID:  req.Id,
		URL: req.Url,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SuccessResponse{
		Message: "Ok",
	}, nil
}

func main() {

	// add a listener address
	lis, err := net.Listen("tcp", fmt.Sprintf(common.DefaultBaseAddressWOPort, os.Getenv("PORT")))
	if err != nil {
		log.Fatalf("ERROR STARTING THE SERVER : %v", err)
	}
	// start the grpc server
	grpcServer := grpc.NewServer()
	db, err := gorm.Open(postgres.Open(db.NewConnectionString(db.PostgresOptions{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     "postgres",
		Password: "admin",
		DBName:   "franchises_db",
		Port:     os.Getenv("POSTGRES_PORT"),
	})), new(gorm.Config))
	if err != nil {
		log.Fatalf("grpc.main: error creating db connection (%v)", err)
	}
	pb.RegisterFranchiseServiceServer(grpcServer, NewFranchiseGRPCServer(
		usecase.NewFranchiseUseCase(
			http.NewScrapFranchiseRepository(scrapToolURL),
			psql.NewCountryRepository(db),
			psql.NewCityRepository(db),
			psql.NewCompanyRepository(db),
			psql.NewLocationRepository(db),
			psql.NewFranchiseRepository(db),
			nil,
		),
	))
	// start serving to the address
	log.Fatal(grpcServer.Serve(lis))
}
