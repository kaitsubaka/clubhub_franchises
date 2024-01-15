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

type UpdateFranchiseDTO struct {
	ID         string
	Title      *string
	SiteName   *string
	LocationID *uint
	Location   *UpdateLocationDTO
}

type UpdateLocationDTO struct {
	Address *string
	ZipCode *string
	City    *string
	Country *string
}

type UpdatedFranchiseDTO struct {
	ID       string
	Title    string
	SiteName string
	URL      string
	Location UpdatedLocationDTO
}

type UpdatedLocationDTO struct {
	Address string
	ZipCode string
	City    string
	Country string
}
