package usecase

import (
	"github.com/google/uuid"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	httprport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/repository/http"
	psqlrport "github.com/kaitsubaka/clubhub_franchises/internal/core/port/repository/psql"
)

type FranchiseUseCase struct {
	franchiseScrapper   httprport.ScrapFranchiseRepository
	countryRepository   psqlrport.CountryRepository
	cityRepository      psqlrport.CityRepository
	companyRepository   psqlrport.CompanyRepository
	locationRepository  psqlrport.LocationRepository
	franchiseRepository psqlrport.FranchiseRepository
}

func NewFranchiseUseCase(
	franchiseScrapper httprport.ScrapFranchiseRepository,
	countryRepository psqlrport.CountryRepository,
	cityRepository psqlrport.CityRepository,
	companyRepository psqlrport.CompanyRepository,
	locationRepository psqlrport.LocationRepository,
	franchiseRepository psqlrport.FranchiseRepository,
) *FranchiseUseCase {
	return &FranchiseUseCase{
		franchiseScrapper:   franchiseScrapper,
		countryRepository:   countryRepository,
		cityRepository:      cityRepository,
		companyRepository:   companyRepository,
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
		Name: func() string {
			if scrapDTO.WhoisData.Administrative != nil {
				return scrapDTO.WhoisData.Administrative.Country
			}
			return ""
		}(),
	})
	if err != nil {
		return err
	}

	cityDTO, err := fuc.cityRepository.Put(dto.CityDTO{
		CountryID: countryDTO.ID,
		Name: func() string {
			if scrapDTO.WhoisData.Administrative != nil {
				return scrapDTO.WhoisData.Administrative.City
			}
			return ""
		}(),
	})
	if err != nil {
		return err
	}

	locationDTO, err := fuc.locationRepository.Put(dto.LocationDTO{
		CityID: cityDTO.ID,
		Address: func() string {
			if scrapDTO.WhoisData.Administrative != nil {
				return scrapDTO.WhoisData.Administrative.Street
			}
			return ""
		}(),
		ZipCode: func() string {
			if scrapDTO.WhoisData.Administrative != nil {
				return scrapDTO.WhoisData.Administrative.PostalCode
			}
			return ""
		}(),
	})
	if err != nil {
		return err
	}
	companyDTO, err := fuc.companyRepository.Put(dto.CompanyDTO{
		ID: uuid.NewString(),
		//TODO: get owner id from onwer created in db
		OwnerID:    uuid.NewString(),
		Name:       scrapDTO.WhoisData.Registrar.Name,
		TaxNumber:  "todo",
		LocationID: locationDTO.ID,
	})
	if err != nil {
		return err
	}

	if _, err := fuc.franchiseRepository.Put(dto.FranchiseDTO{
		ID:                   f.ID,
		CompanyID:            companyDTO.ID,
		Title:                scrapDTO.HTMLData.Title,
		SiteName:             scrapDTO.HTMLData.SiteName,
		Description:          scrapDTO.HTMLData.Description,
		Image:                scrapDTO.HTMLData.Image,
		URL:                  f.URL,
		Protocol:             scrapDTO.Protocol,
		DomainJumps:          scrapDTO.Jumps,
		ServerNames:          scrapDTO.WhoisData.Domain.NameServers,
		DomainCreationDate:   scrapDTO.WhoisData.Domain.CreatedDate,
		DomainExpirationDate: scrapDTO.WhoisData.Domain.ExpirationDate,
		RegistrantName:       scrapDTO.WhoisData.Registrant.Name,
		ContactEmail:         scrapDTO.WhoisData.Registrant.Email,
		LocationID:           locationDTO.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (fuc *FranchiseUseCase) Update(u dto.UpdateFranchiseDTO) (dto.UpdatedFranchiseDTO, error) {
	updateLocDTO := dto.UpdatedLocationDTO{}
	if u.Location != nil {
		locID, err := fuc.franchiseRepository.ConsultLocationID(u.ID)
		if err != nil {
			return dto.UpdatedFranchiseDTO{}, err
		}
		oldLocationDTO, err := fuc.locationRepository.GetByID(locID)
		if err != nil {
			return dto.UpdatedFranchiseDTO{}, err
		}
		oldCityDTO, err := fuc.cityRepository.GetByID(oldLocationDTO.CityID)
		if err != nil {
			return dto.UpdatedFranchiseDTO{}, err
		}
		oldCountryDTO, err := fuc.countryRepository.GetByID(oldCityDTO.CountryID)
		if err != nil {
			return dto.UpdatedFranchiseDTO{}, err
		}

		countryDTO, err := fuc.countryRepository.Put(dto.CountryDTO{
			Name: func() string {
				if u.Location.Country != nil {
					return *u.Location.Country
				}
				return oldCountryDTO.Name
			}(),
		})
		if err != nil {
			return dto.UpdatedFranchiseDTO{}, err
		}

		cityDTO, err := fuc.cityRepository.Put(dto.CityDTO{
			CountryID: countryDTO.ID,
			Name: func() string {
				if u.Location.City != nil {
					return *u.Location.City
				}
				return oldCityDTO.Name
			}(),
		})
		if err != nil {
			return dto.UpdatedFranchiseDTO{}, err
		}

		locationDTO, err := fuc.locationRepository.Put(dto.LocationDTO{
			CityID: cityDTO.ID,
			Address: func() string {
				if u.Location.Address != nil {
					return *u.Location.Address
				}
				return oldLocationDTO.Address
			}(),
			ZipCode: func() string {
				if u.Location.ZipCode != nil {
					return *u.Location.ZipCode
				}
				return oldLocationDTO.ZipCode
			}(),
		})
		if err != nil {
			return dto.UpdatedFranchiseDTO{}, err
		}

		u.LocationID = &locationDTO.ID
		updateLocDTO.Address = locationDTO.Address
		updateLocDTO.City = cityDTO.Name
		updateLocDTO.ZipCode = locationDTO.ZipCode
		updateLocDTO.Country = countryDTO.Name
	}
	updatedDTO, err := fuc.franchiseRepository.Update(u)
	if err != nil {
		return dto.UpdatedFranchiseDTO{}, err
	}
	return dto.UpdatedFranchiseDTO{
		ID:       updatedDTO.ID,
		Title:    updatedDTO.Title,
		SiteName: updatedDTO.SiteName,
		URL:      updatedDTO.URL,
		Location: updateLocDTO,
	}, nil
}
