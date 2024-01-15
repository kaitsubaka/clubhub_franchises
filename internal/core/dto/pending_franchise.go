package dto

type PendingFranchiseDTO struct {
	ID     string
	URL    string
	Status string
	Error  *string
}

type CreateNewFranchizeDTO struct {
	URL string
}
