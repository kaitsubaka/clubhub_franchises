package psql

import (
	"errors"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	pubdto "github.com/kaitsubaka/clubhub_franchises/internal/infra/dto"
	"gorm.io/gorm"
)

type FranchiseRepository struct {
	db *gorm.DB
}

func NewFranchiseRepository(db *gorm.DB) *FranchiseRepository {
	return &FranchiseRepository{db: db}
}

func (fr *FranchiseRepository) Put(f dto.FranchiseDTO) (dto.FranchiseDTO, error) {

	model := pubdto.FranchiseModel{}

	trx := fr.db.First(&model, "title = ? AND company_id = ?",
		f.Title,
		f.CompanyID,
	)

	if errors.Is(trx.Error, gorm.ErrRecordNotFound) {
		trx = fr.db.Create(&pubdto.FranchiseModel{
			ID:                   f.ID,
			CompanyID:            f.CompanyID,
			Title:                f.Title,
			SiteName:             f.SiteName,
			Description:          f.Description,
			Image:                f.Image,
			URL:                  f.URL,
			Protocol:             f.Protocol,
			DomainJumps:          f.DomainJumps,
			ServerNames:          f.ServerNames,
			DomainCreationDate:   f.DomainCreationDate,
			DomainExpirationDate: f.DomainExpirationDate,
			RegistrantName:       f.RegistrantName,
			ContactEmail:         f.ContactEmail,
			LocationID:           f.LocationID,
		})

		return f, trx.Error
	}

	if trx.Error != nil {
		return dto.FranchiseDTO{}, trx.Error
	}

	return f, nil
}
