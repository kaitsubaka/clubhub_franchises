package psql

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type CompanyRepository interface {
	Put(c dto.CompanyDTO) (dto.CompanyDTO, error)
}
