package scraper

import (
	"net/url"

	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-multierror"
)

type Scraper struct {
	err error
}

func NewScraper() *Scraper {
	return &Scraper{}
}

func (s *Scraper) HandleError() colly.ErrorCallback {
	return func(_ *colly.Response, err error) {
		s.err = multierror.Append(s.err, err)
	}
}

func (s *Scraper) SaveHTML() colly.ResponseCallback {
	return func(r *colly.Response) {
		if err := r.Save(urlToFilename(r.Request.URL)); err != nil {
			s.err = multierror.Append(s.err, err)
		}
	}
}

func (s *Scraper) Err() error {
	return s.err
}

func urlToFilename(u *url.URL) string {
	return u.Host + `.html`
}
