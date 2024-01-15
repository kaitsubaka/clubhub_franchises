package psql

import "github.com/kaitsubaka/clubhub_franchises/internal/core/dto"

type LocationRepository interface {
	Put(c dto.LocationDTO) (dto.LocationDTO, error)
	GetByID(id uint) (dto.LocationDTO, error)
}
