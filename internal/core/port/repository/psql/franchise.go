package psql

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type FranchiseRepository interface {
	Put(c dto.FranchiseDTO) (dto.FranchiseDTO, error)
}
