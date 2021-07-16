package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gocolly/colly/v2"
	"github.com/shihanng/webdl/scraper"
)

func main() {
	logger := &log.Logger{
		Level:   log.ErrorLevel,
		Handler: cli.New(os.Stderr),
	}

	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Println(`Archive web pages to disk.  When <url1> is www.google.com,
this tool will download the page and save it as www.google.com.html.

Usage:
  webdl [options] <url1> <url2> ...

Options:`)
		fs.PrintDefaults()
	}

	debug := fs.Bool("debug", false, "show debug log")
	metadata := fs.Bool("metadata", false, "show metadata")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.WithError(err).Fatal("failed to parse flags")
	}

	if *debug {
		logger.Level = log.DebugLevel
	}

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	s := scraper.NewScraper(logger)

	c := colly.NewCollector()

	c.OnRequest(s.VisitLog())
	c.OnResponse(s.SaveHTML())
	if *metadata {
		c.OnHTML("img[src]", s.CountImage())
		c.OnHTML("a[href]", s.CountLink())
	}
	c.OnError(s.HandleError())

	// Visits all URLs
	for _, url := range args {
		if err := c.Visit(url); err != nil {
			logger.WithError(err).Fatal("failed to visit")
		}
	}

	if err := s.Err(); err != nil {
		logger.WithError(err).Fatal("error occurred during scrapping")
	}

	if *metadata {
		for _, p := range s.Pages {
			fmt.Printf("%v\n\n", p)
		}
	}
}
