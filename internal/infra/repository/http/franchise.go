package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	liburl "net/url"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/promise"
	pubdto "github.com/kaitsubaka/clubhub_franchises/internal/infra/dto"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"go.uber.org/multierr"
	"golang.org/x/net/html"
)

type ScrapFranchiseRepository struct {
	toolURL string
}

func NewScrapFranchiseRepository(toolURL string) *ScrapFranchiseRepository {
	return &ScrapFranchiseRepository{
		toolURL: toolURL,
	}
}

func (sfr *ScrapFranchiseRepository) Scrap(url string) (scrap pubdto.ScrapDTO, errors error) {
	done := make(chan struct{})
	defer close(done)
	whoisData := sfr.parseWhoisData(url)
	pnjData := sfr.parseProtocolAndJumps(url)
	htmlData := sfr.parseHTMLData(url)
	resolved := promise.All(done, htmlData, pnjData, whoisData)
	var localErrors error
	for v := range resolved {
		if err, ok := v.(error); ok {
			localErrors = multierr.Append(errors, err)
			continue
		}
		switch v.(type) {
		case pubdto.HTMLDataDTO:
			scrap.HTMLData = v.(pubdto.HTMLDataDTO)
		case pubdto.ProtocolAndJumpsDTO:
			scrap.ProtocolAndJumpsDTO = v.(pubdto.ProtocolAndJumpsDTO)
		case whoisparser.WhoisInfo:
			scrap.WhoisData = v.(whoisparser.WhoisInfo)
		}
	}
	if localErrors != nil {
		log.Println("[INFO] Scrap: ", localErrors)
	}
	return
}

func (sfr *ScrapFranchiseRepository) parseWhoisData(url string) <-chan any {
	prom := make(chan any)
	go func() {
		defer close(prom)

		r, err := liburl.Parse(url)
		if err != nil {
			prom <- err
			return
		}

		whois_raw, err := whois.Whois(r.Host)
		if err != nil {
			prom <- err
			return
		}

		result, err := whoisparser.Parse(whois_raw)

		if err != nil {
			prom <- err
			return
		}

		prom <- result
	}()
	return prom
}

func (sfr *ScrapFranchiseRepository) parseProtocolAndJumps(url string) <-chan any {
	prom := make(chan any)

	go func() {
		defer close(prom)

		resp, err := http.Get(fmt.Sprintf(sfr.toolURL, url))
		if err != nil {
			prom <- err
			return
		}

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			prom <- err
			return
		}

		sslLabsResponse := new(pubdto.SSLLabsResponseDTO)

		err = json.Unmarshal(b, sslLabsResponse)
		if err != nil {
			prom <- err
			return
		}

		prom <- pubdto.ProtocolAndJumpsDTO{
			Protocol: sslLabsResponse.Protocol,
			Jumps:    len(sslLabsResponse.Endpoints),
		}
	}()

	return prom
}

func (sfr *ScrapFranchiseRepository) parseHTMLData(url string) <-chan any {
	prom := make(chan any)
	go func() {
		defer close(prom)
		resp, err := http.Get(url)
		if err != nil {
			prom <- err
			return
		}
		defer resp.Body.Close()
		data := sfr.extractHeaderData(resp.Body)
		prom <- data
	}()
	return prom
}

func (sfr *ScrapFranchiseRepository) extractHeaderData(b io.Reader) pubdto.HTMLDataDTO {
	htmlTokenizer := html.NewTokenizer(b)
	titleFound := false
	data := new(pubdto.HTMLDataDTO)

	for {
		tk := htmlTokenizer.Next()
		switch tk {
		case html.ErrorToken:
			return *data
		case html.StartTagToken, html.SelfClosingTagToken:
			t := htmlTokenizer.Token()
			if t.Data == `body` {
				return *data
			}
			if t.Data == "title" {
				titleFound = true
			}
			if t.Data == "meta" {
				desc, ok := sfr.parseMetaProperty("description", t)
				if ok {
					data.Description = desc
				}

				ogTitle, ok := sfr.parseMetaProperty("og:title", t)
				if ok {
					data.Title = ogTitle
				}

				ogDesc, ok := sfr.parseMetaProperty("og:description", t)
				if ok {
					data.Description = ogDesc
				}

				ogImage, ok := sfr.parseMetaProperty("og:image", t)
				if ok {
					data.Image = ogImage
				}

				ogSiteName, ok := sfr.parseMetaProperty("og:site_name", t)
				if ok {
					data.SiteName = ogSiteName
				}
			}
		case html.TextToken:
			if titleFound {
				t := htmlTokenizer.Token()
				data.Title = t.Data
				titleFound = false
			}
		}
	}
}

func (sfr *ScrapFranchiseRepository) parseMetaProperty(prop string, t html.Token) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}
	return
}
