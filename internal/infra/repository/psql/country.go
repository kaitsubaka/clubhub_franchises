package psql

import (
	"errors"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	pubdto "github.com/kaitsubaka/clubhub_franchises/internal/infra/dto"
	"gorm.io/gorm"
)

type CountryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (lr *CountryRepository) Put(c dto.CountryDTO) (dto.CountryDTO, error) {
	localCountry := new(pubdto.CountryModel)
	trx := lr.db.First(localCountry, "name = ?",
		c.Name,
	)
	if trx.Error != nil {
		if errors.Is(trx.Error, gorm.ErrRecordNotFound) {
			country := &pubdto.CountryModel{
				Name: c.Name,
			}
			trx = lr.db.Create(country)
			c.ID = country.ID
			return c, trx.Error
		}
		return dto.CountryDTO{}, trx.Error
	}
	c.ID = localCountry.ID
	return c, nil
}
