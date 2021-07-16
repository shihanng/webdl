package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/shihanng/webdl/scraper"
)

const (
	nConsumer    = 2
	maxQueueSize = 10000
)

func main() {
	logger := &log.Logger{
		Level:   log.ErrorLevel,
		Handler: cli.New(os.Stderr),
	}

	fs := flag.NewFlagSet("", flag.ExitOnError)
	debug := fs.Bool("debug", false, "show debug log")
	metadata := fs.Bool("metadata", false, "show metadata")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.WithError(err).Fatal("failed to parse flags")
	}

	if *debug {
		logger.Level = log.DebugLevel
	}

	args := fs.Args()

	s := scraper.NewScraper(logger)

	c := colly.NewCollector()

	c.OnRequest(s.VisitLog())
	c.OnResponse(s.SaveHTML())
	if *metadata {
		c.OnHTML("img[src]", s.CountImage())
		c.OnHTML("a[href]", s.CountLink())
	}
	c.OnError(s.HandleError())

	// Create request queue.
	q, err := queue.New(nConsumer, &queue.InMemoryQueueStorage{MaxSize: maxQueueSize})
	if err != nil {
		logger.WithError(err).Fatal("failed to create queue")
	}

	// Enqueue URLs to visit.
	for _, url := range args {
		if err := q.AddURL(url); err != nil {
			logger.WithError(err).Fatal("failed to enqueue")
		}
	}

	if err := q.Run(c); err != nil {
		logger.WithError(err).Fatal("failed to consume queue")
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
