package psql

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type CountryRepository interface {
	Put(c dto.CountryDTO) (dto.CountryDTO, error)
}
