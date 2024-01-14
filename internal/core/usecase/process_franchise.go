package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/queue/event"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	pb "github.com/kaitsubaka/clubhub_franchises/proto"
	"google.golang.org/grpc"
)

type ProcessFranchiseUseCase struct {
	franchiseServiceClient pb.FranchiseServiceClient
}

func NewProcessFranchiseUseCase(serviceConn grpc.ClientConnInterface) *ProcessFranchiseUseCase {
	return &ProcessFranchiseUseCase{
		franchiseServiceClient: pb.NewFranchiseServiceClient(serviceConn),
	}
}

func (pfuc *ProcessFranchiseUseCase) Subscribe(e event.Event) error {
	fDTO, ok := e.Data.(dto.PendingFranchiseDTO)
	if !ok {
		return errors.New("ProcessFranchiseUseCase.Subscribe:the event data cannot be handled")
	}
	return pfuc.process(fDTO)
}

func (pfuc *ProcessFranchiseUseCase) process(in dto.PendingFranchiseDTO) error {
	res, err := pfuc.franchiseServiceClient.Create(context.Background(), &pb.CreateFranchiseRequest{
		Id:  in.ID,
		Url: in.URL,
	})
	if err != nil {
		return err
	}
	log.Println("[INFO] process: ", res.Message)
	return nil
}
