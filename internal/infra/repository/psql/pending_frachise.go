package psql

import (
	"errors"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	publicdto "github.com/kaitsubaka/clubhub_franchises/internal/infra/dto"
	"gorm.io/gorm"
)

type PendingFranchizeRepository struct {
	db *gorm.DB
}

func NewPendingFranchizeRepository(db *gorm.DB) *PendingFranchizeRepository {
	return &PendingFranchizeRepository{
		db: db,
	}
}

func (pfr *PendingFranchizeRepository) Put(n dto.PendingFranchiseDTO) (dto.PendingFranchiseDTO, error) {
	localNewFranchise := new(publicdto.PendingFranchizeModel)
	trx := pfr.db.First(localNewFranchise, "url = ?", n.URL)
	if trx.Error != nil {
		if errors.Is(trx.Error, gorm.ErrRecordNotFound) {
			trx = pfr.db.Create(&publicdto.PendingFranchizeModel{
				ID:     n.ID,
				URL:    n.URL,
				Status: "CREATED",
			})
			return n, trx.Error
		}
		return dto.PendingFranchiseDTO{}, trx.Error
	}
	n.ID = localNewFranchise.ID
	return n, nil
}

func (pfr *PendingFranchizeRepository) UpdateStatus(n dto.PendingFranchiseDTO) error {
	trx := pfr.db.Save(&publicdto.PendingFranchizeModel{
		ID:     n.ID,
		Status: n.Status,
		Error:  n.Error,
	})
	if trx.Error != nil {
		return trx.Error
	}
	return nil
}

func (pfr *PendingFranchizeRepository) GetStatusByID(id string) (string, error) {
	pm := new(publicdto.PendingFranchizeModel)
	trx := pfr.db.First(pm, "id = ?", id)
	if trx.Error != nil {
		return "", trx.Error
	}
	return pm.Status, nil
}
