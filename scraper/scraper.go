package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
		site := urlToFilename(r.Request.URL)

		s.Pages[site] = &Page{
			Site:      filepath.Join(r.Request.URL.Host, r.Request.URL.Path),
			LastFetch: time.Now(),
		}

		// Create directory for specific domain.
		if err := os.MkdirAll(r.Request.URL.Host, os.ModePerm); err != nil {
			s.err = multierror.Append(s.err, err)
			return
		}

		// Save the HTML content.
		if err := r.Save(urlToFilename(r.Request.URL)); err != nil {
			s.err = multierror.Append(s.err, err)
		}
	}
}

func (s *Scraper) SaveAsset(attr string) colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		src := e.Attr(attr)

		// Skipping asset that has full URL
		if src == "" || strings.HasPrefix(src, "http") {
			s.log.WithField(attr, src).Debug("skip downloading")
			return
		}

		// Construct filename of the asset
		target := filepath.Join(e.Request.URL.Host, src)

		if err := dowloadFile(e.Request.AbsoluteURL(src), target); err != nil {
			s.err = multierror.Append(s.err, err)
			return
		}
	}
}

func (s *Scraper) CountImage() colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		site := urlToFilename(e.Request.URL)
		if _, ok := s.Pages[site]; ok {
			s.Pages[site].NumImages += 1
		}
	}
}

func (s *Scraper) CountLink() colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		site := urlToFilename(e.Request.URL)
		if _, ok := s.Pages[site]; ok {
			s.Pages[site].NumLinks += 1
		}
	}
}

func (s *Scraper) Err() error {
	return s.err
}

func urlToFilename(u *url.URL) string {
	filename := "index"

	path := strings.TrimSuffix(u.Path, filepath.Ext(u.Path))
	path = strings.TrimPrefix(path, "/")

	if path != "" {
		filename = strings.Replace(path, "/", "_", -1)
	}

	return filepath.Join(u.Host, filename+`.html`)
}

func dowloadFile(url string, target string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
