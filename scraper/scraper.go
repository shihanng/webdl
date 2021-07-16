package scraper

import (
	"net/url"

	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-multierror"
)

type Scraper struct {
	log *log.Logger
	err error
}

func NewScraper(logger *log.Logger) *Scraper {
	return &Scraper{log: logger}
}

func (s *Scraper) VisitLog() colly.RequestCallback {
	return func(r *colly.Request) {
		s.log.WithField("url", r.URL.String()).Debug("visiting")
	}
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
