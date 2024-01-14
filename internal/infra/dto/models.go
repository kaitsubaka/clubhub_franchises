package dto

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PendingFranchizeModel struct {
	gorm.Model
	ID  string
	URL string
}

func (PendingFranchizeModel) TableName() string {
	return "pending_franchises"
}

type CompanyModel struct {
	gorm.Model
	ID             string `gorm:"primarykey"`
	Name           string
	CompanyOwnerID string
	TaxNumber      string
	LocationID     string
	Location       LocationModel `gorm:"foreignKey:LocationID;references:ID"`
}

func (CompanyModel) TableName() string {
	return "company"
}

type FranchiseModel struct {
	gorm.Model
	ID                   string `gorm:"primarykey"`
	CompanyID            string
	Title                string
	SiteName             string
	Description          string
	Image                string
	URL                  string
	Protocol             string
	DomainJumps          int
	ServerNames          pq.StringArray `gorm:"type:text[]"`
	DomainCreationDate   string
	DomainExpirationDate string
	RegistrantName       string
	ContactEmail         string
	LocationID           uint
	Company              CompanyModel  `gorm:"foreignKey:CompanyID;references:ID"`
	Location             LocationModel `gorm:"foreignKey:LocationID;references:ID"`
}

func (FranchiseModel) TableName() string {
	return "franchises"
}

type LocationModel struct {
	gorm.Model
	Address string
	ZipCode string
	CityID  uint
}

func (LocationModel) TableName() string {
	return "locations"
}
