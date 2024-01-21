package psql

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/db/test"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SuiteFranchiseRepository struct {
	suite.Suite
	db *gorm.DB
}

func TestSuiteFranchiseRepository(t *testing.T) {
	suite.Run(t, new(SuiteFranchiseRepository))
}

func (s *SuiteFranchiseRepository) SetupSuite() {
	s.db = test.NewPSQLTestDB(filepath.Join("../../../../scripts/db/test", "franchise.sql"), s.T())
}

func (s *SuiteFranchiseRepository) TestPut() {

	fr := NewFranchiseRepository(s.db)
	s.Run("it should return the new franchise created if the franchise is not found in the db", func() {
		d := time.Now().Format(time.RFC1123Z)
		id := uuid.NewString()
		got, err := fr.Put(dto.FranchiseDTO{
			ID:                   id,
			CompanyID:            id,
			ServerNames:          []string{"test"},
			DomainCreationDate:   d,
			DomainExpirationDate: d,
		})
		s.NoError(err)
		s.Equal(
			dto.FranchiseDTO{
				ID: id, CompanyID: id, Title: "", SiteName: "", Description: "", Image: "", URL: "", Protocol: "", DomainJumps: 0, ServerNames: []string{"test"}, DomainCreationDate: d, DomainExpirationDate: d, RegistrantName: "", ContactEmail: "", LocationID: 0x0,
			},
			got,
		)
	})
}

func (s *SuiteFranchiseRepository) TestUpdate() {
	fr := NewFranchiseRepository(s.db)
	s.Run("it should return the updated franchise", func() {
		d := time.Now().Format(time.RFC1123Z)

		id := uuid.NewString()
		_, err := fr.Put(dto.FranchiseDTO{
			ID:                   id,
			CompanyID:            id,
			ServerNames:          []string{"test"},
			DomainCreationDate:   d,
			DomainExpirationDate: d,
		})
		s.NoError(err)
		title := "updated"
		got, err := fr.Update(dto.UpdateFranchiseDTO{
			ID:    id,
			Title: &title,
		})
		s.NoError(err)
		s.Equal(
			dto.FranchiseDTO{
				ID: id, Title: title,
			},
			got,
		)
	})
}

func (s *SuiteFranchiseRepository) TestConsultLocationID() {
	fr := NewFranchiseRepository(s.db)
	s.Run("it should return the location id from the franchise consulted by id", func() {
		d := time.Now().Format(time.RFC1123Z)
		id := uuid.NewString()
		locID := uint(1)
		_, err := fr.Put(dto.FranchiseDTO{
			ID:                   id,
			CompanyID:            id,
			ServerNames:          []string{"test"},
			DomainCreationDate:   d,
			DomainExpirationDate: d,
			LocationID:           locID,
		})
		s.NoError(err)
		got, err := fr.ConsultLocationID(id)
		s.NoError(err)
		s.Equal(
			locID,
			got,
		)
	})
}
