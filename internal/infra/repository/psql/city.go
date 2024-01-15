package psql

import (
	"errors"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	pubdto "github.com/kaitsubaka/clubhub_franchises/internal/infra/dto"
	"gorm.io/gorm"
)

type CityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{db: db}
}

func (lr *CityRepository) Put(c dto.CityDTO) (dto.CityDTO, error) {
	localCity := new(pubdto.CityModel)
	trx := lr.db.First(localCity, "country_id = ? AND name = ?",
		c.CountryID,
		c.Name,
	)
	if trx.Error != nil {
		if errors.Is(trx.Error, gorm.ErrRecordNotFound) {
			city := &pubdto.CityModel{
				CountryID: c.CountryID,
				Name:      c.Name,
			}
			trx = lr.db.Create(city)
			c.ID = city.ID
			return c, trx.Error
		}
		return dto.CityDTO{}, trx.Error
	}
	c.ID = localCity.ID
	return c, nil
}

func (lr *CityRepository) GetByID(id uint) (dto.CityDTO, error) {
	localModel := new(pubdto.CityModel)
	trx := lr.db.Find(localModel, id)
	if trx.Error != nil {
		return dto.CityDTO{}, trx.Error
	}
	return dto.CityDTO{
		ID:        localModel.ID,
		CountryID: localModel.CountryID,
		Name:      localModel.Name,
	}, nil
}
