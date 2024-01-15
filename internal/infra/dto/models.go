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
	ID         string `gorm:"primarykey"`
	Name       string
	OwnerID    string
	TaxNumber  string
	LocationID uint
	Location   LocationModel `gorm:"foreignKey:LocationID;references:ID"`
}

func (CompanyModel) TableName() string {
	return "companies"
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

type CityModel struct {
	gorm.Model
	CountryID uint
	Name      string
}

func (CityModel) TableName() string {
	return "cities"
}

type CountryModel struct {
	gorm.Model
	Name string
}

func (CountryModel) TableName() string {
	return "countries"
}
