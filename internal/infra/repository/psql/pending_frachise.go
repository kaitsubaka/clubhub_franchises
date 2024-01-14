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
				ID:  n.ID,
				URL: n.URL,
			})
			return n, trx.Error
		}
		return dto.PendingFranchiseDTO{}, trx.Error
	}
	n.ID = localNewFranchise.ID
	return n, nil
}
