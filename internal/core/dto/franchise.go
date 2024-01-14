package dto

type FranchiseDTO struct {
	ID                   string
	CompanyID            string
	Title                string
	SiteName             string
	Description          string
	Image                string
	URL                  string
	Protocol             string
	DomainJumps          int
	ServerNames          []string
	DomainCreationDate   string
	DomainExpirationDate string
	RegistrantName       string
	ContactEmail         string
	LocationID           uint
}
