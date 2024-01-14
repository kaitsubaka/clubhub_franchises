package http

import pubdto "github.com/kaitsubaka/clubhub_franchises/internal/infra/dto"

type ScrapFranchiseRepository interface {
	Scrap(url string) (scrap pubdto.ScrapDTO, errors error)
}
