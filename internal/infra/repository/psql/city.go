package psql

import (
	"errors"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	"gorm.io/gorm"
)

type CityModel struct {
	gorm.Model
	CountryID uint
	Name      string
}

func (CityModel) TableName() string {
	return "cities"
}

type CityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{db: db}
}

func (lr *CityRepository) Put(c dto.CityDTO) (dto.CityDTO, error) {
	localCity := new(CityModel)
	trx := lr.db.First(localCity, "country_id = ? AND name = ?",
		c.CountryID,
		c.Name,
	)
	if trx.Error != nil {
		if errors.Is(trx.Error, gorm.ErrRecordNotFound) {
			city := &CityModel{
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
