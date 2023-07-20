package goproxyscrape

type Scraper interface {
	Scrape() error
}
