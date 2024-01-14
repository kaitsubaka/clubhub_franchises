package usecase

import (
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	httprport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/repository/http"
	psqlrport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/repository/psql"
)

type FranchiseUseCase struct {
	franchiseScrapper   httprport.ScrapFranchiseRepository
	countryRepository   psqlrport.CountryRepository
	cityRepository      psqlrport.CityRepository
	locationRepository  psqlrport.LocationRepository
	franchiseRepository psqlrport.FranchiseRepository
}

func NewFranchiseUseCase(
	franchiseScrapper httprport.ScrapFranchiseRepository,
	countryRepository psqlrport.CountryRepository,
	locationRepository psqlrport.LocationRepository,
	franchiseRepository psqlrport.FranchiseRepository,
) *FranchiseUseCase {
	return &FranchiseUseCase{
		franchiseScrapper:   franchiseScrapper,
		countryRepository:   countryRepository,
		locationRepository:  locationRepository,
		franchiseRepository: franchiseRepository,
	}
}

func (fuc *FranchiseUseCase) Create(f dto.PendingFranchiseDTO) error {
	scrapDTO, err := fuc.franchiseScrapper.Scrap(f.URL)
	if err != nil {
		return err
	}
	countryDTO, err := fuc.countryRepository.Put(dto.CountryDTO{
		Name: scrapDTO.WhoisData.Administrative.Country,
	})
	if err != nil {
		return err
	}
	cityDTO, err := fuc.cityRepository.Put(dto.CityDTO{
		CountryID: countryDTO.ID,
		Name:      scrapDTO.WhoisData.Administrative.City,
	})
	if err != nil {
		return err
	}
	locationDTO, err := fuc.locationRepository.Put(dto.LocationDTO{
		CityID:  cityDTO.ID,
		Address: scrapDTO.WhoisData.Administrative.Street,
		ZipCode: scrapDTO.WhoisData.Administrative.PostalCode,
	})
	if err != nil {
		return err
	}

	if _, err := fuc.franchiseRepository.Put(dto.FranchiseDTO{
		LocationID: locationDTO.ID,
	}); err != nil {
		return err
	}

	return nil
}
