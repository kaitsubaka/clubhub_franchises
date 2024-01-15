package psql

import (
	"errors"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	pubdto "github.com/kaitsubaka/clubhub_franchises/internal/infra/dto"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (cr *CompanyRepository) Put(c dto.CompanyDTO) (dto.CompanyDTO, error) {
	localCompany := new(pubdto.CompanyModel)
	trx := cr.db.First(localCompany, "name = ?", c.Name)
	if trx.Error != nil {
		if errors.Is(trx.Error, gorm.ErrRecordNotFound) {
			company := &pubdto.CompanyModel{
				ID:         c.ID,
				OwnerID:    c.OwnerID,
				Name:       c.Name,
				LocationID: c.LocationID,
			}
			trx = cr.db.Create(company)
			c.ID = company.ID
			return c, trx.Error
		}
		return dto.CompanyDTO{}, trx.Error
	}
	c.ID = localCompany.ID
	return c, nil
}
