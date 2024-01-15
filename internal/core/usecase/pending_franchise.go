package usecase

import (
	"github.com/google/uuid"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/queue/event"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	qport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/queue"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/port/repository/psql"
)

type PendingFranchiseUseCase struct {
	newFranchizeRepository     psql.PendingFranchiseRepository
	franchiseCreationPublisher qport.Publisher
}

func NewPendingFranchiseUseCase(newFranchizeRepository psql.PendingFranchiseRepository, franchiseCreationPublisher qport.Publisher) *PendingFranchiseUseCase {
	return &PendingFranchiseUseCase{
		newFranchizeRepository:     newFranchizeRepository,
		franchiseCreationPublisher: franchiseCreationPublisher,
	}
}

func (nfs *PendingFranchiseUseCase) Create(in dto.CreateNewFranchizeDTO) (dto.PendingFranchiseDTO, error) {
	createdFranchise, err := nfs.newFranchizeRepository.Put(dto.PendingFranchiseDTO{
		ID:  uuid.NewString(),
		URL: in.URL,
	})
	if err != nil {
		return dto.PendingFranchiseDTO{}, err
	}

	if err := nfs.franchiseCreationPublisher.Publish(event.Event{
		ID:   createdFranchise.ID,
		Data: createdFranchise,
	}); err != nil {
		return dto.PendingFranchiseDTO{}, err
	}

	return createdFranchise, nil
}
