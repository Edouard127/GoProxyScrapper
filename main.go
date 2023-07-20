package goproxyscrape

var Dialer = NewHTTPDialer()

func Scrape() []Proxy {
	var proxies []Proxy
	for _, scraper := range Webpages {
		proxies = append(proxies, scraper()...)
	}
	return proxies
}
