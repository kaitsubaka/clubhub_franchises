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

func (fr *FranchiseRepository) Update(f dto.UpdateFranchiseDTO) (dto.FranchiseDTO, error) {
	model := new(pubdto.FranchiseModel)
	trx := fr.db.First(model, "id = ?", f.ID)
	if trx.Error != nil {
		return dto.FranchiseDTO{}, trx.Error
	}
	if f.Title != nil {
		model.Title = *f.Title
	}
	if f.SiteName != nil {
		model.SiteName = *f.SiteName
	}
	if f.LocationID != nil {
		model.LocationID = *f.LocationID
	}
	trx = fr.db.Save(model)
	if trx.Error != nil {
		return dto.FranchiseDTO{}, trx.Error
	}
	return dto.FranchiseDTO{
		ID:       model.ID,
		SiteName: model.SiteName,
		Title:    model.Title,
		URL:      model.URL,
	}, nil
}

func (fr *FranchiseRepository) ConsultLocationID(id string) (uint, error) {
	model := new(pubdto.FranchiseModel)
	trx := fr.db.First(model, "id = ?", id)
	if trx.Error != nil {
		return 0, trx.Error
	}
	return model.LocationID, nil
}
