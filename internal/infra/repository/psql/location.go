package psql

import (
	"errors"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	publicdto "github.com/kaitsubaka/clubhub_franchises/internal/infra/dto"
	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (lr *LocationRepository) Put(l dto.LocationDTO) (dto.LocationDTO, error) {
	localLocation := new(publicdto.LocationModel)
	trx := lr.db.First(localLocation, "address = ? AND city_id = ? AND zip_code = ?",
		l.Address,
		l.CityID,
		l.ZipCode,
	)
	if trx.Error != nil {
		if errors.Is(trx.Error, gorm.ErrRecordNotFound) {
			location := &publicdto.LocationModel{
				Address: l.Address,
				CityID:  l.CityID,
				ZipCode: l.ZipCode,
			}
			trx = lr.db.Create(location)
			l.ID = location.ID
			return l, trx.Error
		}
		return dto.LocationDTO{}, trx.Error
	}
	l.ID = localLocation.ID
	return l, nil
}

func (lr *LocationRepository) GetByID(id uint) (dto.LocationDTO, error) {
	localLocation := new(publicdto.LocationModel)
	trx := lr.db.Find(localLocation, id)
	if trx.Error != nil {
		return dto.LocationDTO{}, trx.Error
	}
	return dto.LocationDTO{
		ID:      localLocation.ID,
		Address: localLocation.Address,
		ZipCode: localLocation.ZipCode,
		CityID:  localLocation.CityID,
	}, nil
}
