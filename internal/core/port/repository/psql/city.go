package psql

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type CityRepository interface {
	Put(c dto.CityDTO) (dto.CityDTO, error)
	GetByID(id uint) (dto.CityDTO, error)
}
