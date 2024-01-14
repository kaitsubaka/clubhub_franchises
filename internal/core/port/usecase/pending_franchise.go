package usecase

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type PendingFranchiseUseCase interface {
	Create(in dto.CreateNewFranchizeDTO) (dto.PendingFranchiseDTO, error)
}
