package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/queue/event"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	psqlrport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/repository/psql"
	pb "github.com/kaitsubaka/clubhub_franchises/proto"
	"google.golang.org/grpc"
)

type ProcessFranchiseUseCase struct {
	franchiseServiceClient     pb.FranchiseServiceClient
	pendingFranchiseRepository psqlrport.PendingFranchiseRepository
}

func NewProcessFranchiseUseCase(serviceConn grpc.ClientConnInterface, pendingFranchiseRepository psqlrport.PendingFranchiseRepository) *ProcessFranchiseUseCase {
	return &ProcessFranchiseUseCase{
		franchiseServiceClient:     pb.NewFranchiseServiceClient(serviceConn),
		pendingFranchiseRepository: pendingFranchiseRepository,
	}
}

func (pfuc *ProcessFranchiseUseCase) Subscribe(e event.Event) error {
	fDTO, ok := e.Data.(dto.PendingFranchiseDTO)
	if !ok {
		return errors.New("ProcessFranchiseUseCase.Subscribe:the event data cannot be handled")
	}
	if err := pfuc.process(fDTO); err != nil {
		errosMsg := fmt.Sprint(err)
		err = pfuc.pendingFranchiseRepository.UpdateStatus(dto.PendingFranchiseDTO{
			ID:     fDTO.ID,
			Error:  &errosMsg,
			Status: "ERROR",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (pfuc *ProcessFranchiseUseCase) process(in dto.PendingFranchiseDTO) error {
	status, err := pfuc.pendingFranchiseRepository.GetStatusByID(in.ID)
	if err != nil {
		return err
	}
	if (status == "COMPLETED" || status == "PROCESSING") && status != "ERROR" {
		log.Println("[INFO] process: the franchise is already being processed by other routine ignoring request")
		return nil
	}
	err = pfuc.pendingFranchiseRepository.UpdateStatus(dto.PendingFranchiseDTO{
		ID:     in.ID,
		Status: "PROCESSING",
	})
	if err != nil {
		return err
	}
	res, err := pfuc.franchiseServiceClient.Create(context.Background(), &pb.CreateFranchiseRequest{
		Id:  in.ID,
		Url: in.URL,
	})
	if err != nil {
		return err
	}
	log.Println("[INFO] process: ", res.Message)
	err = pfuc.pendingFranchiseRepository.UpdateStatus(dto.PendingFranchiseDTO{
		ID:     in.ID,
		Status: "COMPLETED",
	})
	if err != nil {
		return err
	}
	return nil
}
