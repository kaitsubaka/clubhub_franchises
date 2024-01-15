package usecase

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type FranchiseUseCase interface {
	Create(in dto.PendingFranchiseDTO) error
	Update(u dto.UpdateFranchiseDTO) (dto.UpdatedFranchiseDTO, error)
}
