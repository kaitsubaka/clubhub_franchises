package psql

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type PendingFranchiseRepository interface {
	Put(n dto.PendingFranchiseDTO) (dto.PendingFranchiseDTO, error)
	UpdateStatus(n dto.PendingFranchiseDTO) error
	GetStatusByID(id string) (string, error)
}
