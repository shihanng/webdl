package scraper

import (
	"fmt"
	"net/url"
	"time"

	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-multierror"
)

type Scraper struct {
	Pages map[string]*Page

	log *log.Logger
	err error
}

type Page struct {
	Site      string
	NumLinks  int
	NumImages int
	LastFetch time.Time
}

func (p *Page) String() string {
	return fmt.Sprintf("site: %s\nnum_links: %d\nimages: %d\nlast_fetch: %s",
		p.Site,
		p.NumLinks,
		p.NumImages,
		p.LastFetch.UTC().Format("Mon Jan _2 2006 15:04 MST"),
	)
}

func NewScraper(logger *log.Logger) *Scraper {
	return &Scraper{
		log:   logger,
		Pages: map[string]*Page{},
	}
}

func (s *Scraper) VisitLog() colly.RequestCallback {
	return func(r *colly.Request) {
		site := urlToFilename(r.URL)

		s.Pages[site] = &Page{
			Site:      site,
			LastFetch: time.Now(),
		}

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

func (s *Scraper) CountImage() colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		site := urlToFilename(e.Request.URL)
		s.Pages[site].NumImages += 1
	}
}

func (s *Scraper) CountLink() colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		site := urlToFilename(e.Request.URL)
		s.Pages[site].NumLinks += 1
	}
}

func (s *Scraper) Err() error {
	return s.err
}

func urlToFilename(u *url.URL) string {
	return u.Host + `.html`
}
