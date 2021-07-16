package scraper

import (
	"net/url"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	err error
}

func NewScraper() *Scraper {
	return &Scraper{}
}

func (s *Scraper) SaveHTML() colly.ResponseCallback {
	return func(r *colly.Response) {
		s.err = r.Save(urlToFilename(r.Request.URL))
	}
}

func (s *Scraper) Err() error {
	return s.err
}

func urlToFilename(u *url.URL) string {
	return u.Host + `.html`
}
