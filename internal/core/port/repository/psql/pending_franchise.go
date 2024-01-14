package psql

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type PendingFranchizeRepository interface {
	Put(n dto.PendingFranchiseDTO) (dto.PendingFranchiseDTO, error)
}
