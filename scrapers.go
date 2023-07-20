package goproxyscrape

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
)

var (
	Webpages []func() []Proxy
)

func init() {
	Webpages = append(Webpages, scrapeAdvancedName)
	//Webpages = append(Webpages, anonymouse)
	//Webpages = append(Webpages, freeProxyCZ)
	Webpages = append(Webpages, freeProxyList)
	Webpages = append(Webpages, freeProxyListCC)
}

func scrapeAdvancedName() []Proxy {
	body, err := Dialer.Dial("https://advanced.name/freeproxy")
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(body)

	proxies := make([]Proxy, 0)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}

	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
		host, _ := s.Find("td[data-ip]").Attr("data-ip")
		port, _ := s.Find("td[data-port]").Attr("data-port")
		country := s.Find("td:nth-child(5) > a").Text()
		var protocol string
		s.Find("td:nth-child(4) > a:not(:last-child)").Each(func(i int, s *goquery.Selection) {
			protocol += s.Text() + ","
		})
		//protocol = protocol[:len(protocol)-1]
		var anonymity string
		if l := s.Find("td:nth-child(4) > a").Length(); l >= 1 {
			anonymity = s.Find("td:nth-child(4) > a:last-child").Text()
		} else {
			anonymity = "Unknown"
		}

		var p Protocol
		var a Anonymity

		p.Parse(protocol)
		a.Parse(anonymity)
		h, _ := base64.StdEncoding.DecodeString(host)
		po, _ := base64.StdEncoding.DecodeString(port)
		atoi, _ := strconv.Atoi(string(po))

		proxies = append(proxies, Proxy{
			Address:   string(h) + ":" + string(po),
			Host:      string(h),
			Port:      uint8(atoi),
			Country:   country,
			Anonymity: a,
			Protocol:  p,
		})
	})

	return proxies
}

func anonymouse() []Proxy {
	body, err := Dialer.Dial("https://anonymouse.cz/proxy-list/")
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(body)

	proxies := make([]Proxy, 0)
	doc, err := goquery.NewDocumentFromReader(reader)

	doc.Find("tbody > tr:not(:first-child)").Each(func(i int, s *goquery.Selection) {
		host := doc.Find("td:first-child").Text()
		port := doc.Find("td:nth-child(2)").Text()
		atoi, _ := strconv.Atoi(port)

		proxies = append(proxies, Proxy{
			Address:   host + ":" + port,
			Host:      host,
			Port:      uint8(atoi),
			Country:   "Unknown",
			Anonymity: UnknownAnonymity,
			Protocol:  UnknownProtocol,
		})
	})

	return proxies
}
func freeProxyCZ() []Proxy {
	body, err := Dialer.Dial("http://free-proxy.cz/en/proxylist/main/1")
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(body)

	proxies := make([]Proxy, 0)
	doc, err := goquery.NewDocumentFromReader(reader)

	lastPage := doc.Find("table#proxy_list tbody > tr")
	fmt.Println(lastPage.Length(), lastPage)
	lastPageInt, _ := strconv.Atoi("1")

	for i := 1; i <= lastPageInt; i++ {
		body, err := Dialer.Dial(fmt.Sprintf("http://free-proxy.cz/en/proxylist/main/%d", i))
		if err != nil {
			panic(err)
		}

		reader := bytes.NewReader(body)

		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			panic(err)
		}

		doc.Find("table#proxy_list tbody > tr").Each(func(i int, s *goquery.Selection) {
			ipCell := s.Find("td > script")
			if ipCell.Length() == 0 {
				return
			}

			encodedStr := ipCell.Text()[len(`Base64.decode("`) : len(ipCell.Text())-len(`")`)]

			// Decode the string from Base64
			decodedBytes, _ := base64.StdEncoding.DecodeString(encodedStr)

			host := string(decodedBytes)
			port := s.Find("td > span.fport").Text()
			atoi, _ := strconv.Atoi(port)
			protocolStr := s.Find("td:nth-child(3) > small").Text()
			countryStr := s.Find("td:nth-child(4) > div > a").Text()
			anonymityStr := s.Find("td:nth-child(7) > small").Text()

			var protocol Protocol
			var anonymity Anonymity

			protocol.Parse(protocolStr)
			anonymity.Parse(anonymityStr)

			proxies = append(proxies, Proxy{
				Address:   host + ":" + port,
				Host:      host,
				Port:      uint8(atoi),
				Country:   countryStr,
				Anonymity: anonymity,
				Protocol:  protocol,
			})
		})
	}

	doc.Find("tbody > tr:not(:first-child)").Each(func(i int, s *goquery.Selection) {
		host := doc.Find("td:first-child").Text()
		port := doc.Find("td:nth-child(2)").Text()
		atoi, _ := strconv.Atoi(port)

		protocolStr := s.Find("td:nth-child(3) > small").Text()
		countryStr := s.Find("td:nth-child(4) > div > a").Text()
		anonymityStr := s.Find("td:nth-child(7) > small").Text()

		var protocol Protocol
		var anonymity Anonymity

		protocol.Parse(protocolStr)
		anonymity.Parse(anonymityStr)

		proxies = append(proxies, Proxy{
			Address:   host + ":" + port,
			Host:      host,
			Port:      uint8(atoi),
			Country:   countryStr,
			Anonymity: anonymity,
			Protocol:  protocol,
		})
	})

	return proxies
}

func freeProxyList() []Proxy {
	controller := func(url string) []Proxy {
		body, err := Dialer.Dial(url)
		if err != nil {
			panic(err)
		}

		reader := bytes.NewReader(body)

		proxies := make([]Proxy, 0)
		doc, err := goquery.NewDocumentFromReader(reader)

		doc.Find("#list tbody tr").Each(func(i int, s *goquery.Selection) {
			host := s.Find("td:first-child").Text()
			port := s.Find("td:nth-child(2)").Text()
			atoi, _ := strconv.Atoi(port)

			var protocolStr string
			var anonymityStr string
			if strings.Contains(url, "https://www.socks-proxy.net/") {
				protocolStr = s.Find("td:nth-child(5)").Text()
			} else if strings.Contains(s.Find("td:nth-child(7)").Text(), "no") {
				protocolStr = "http"
			} else {
				protocolStr = "https"
			}

			if strings.Contains(url, "https://www.socks-proxy.net/") {
				anonymityStr = s.Find("td:nth-child(6)").Text()
			} else {
				anonymityStr = s.Find("td:nth-child(5)").Text()
			}

			countryStr := s.Find("td:nth-child(3)").Text()

			var protocol Protocol
			var anonymity Anonymity

			protocol.Parse(protocolStr)
			anonymity.Parse(anonymityStr)

			proxies = append(proxies, Proxy{
				Address:   host + ":" + port,
				Host:      host,
				Port:      uint8(atoi),
				Country:   countryStr,
				Anonymity: anonymity,
				Protocol:  protocol,
			})
		})

		return proxies
	}

	socksProxyNet := controller("https://www.socks-proxy.net/")
	sslProxiesOrg := controller("https://www.sslproxies.org/")
	freeProxyListNet := controller("https://free-proxy-list.net/")
	usProxyOrg := controller("https://us-proxy.org/")
	proxyNovaCom := controller("https://www.proxynova.com/proxy-server-list/")
	proxyListOrg := controller("https://proxy-list.org/english/index.php?p=1")
	proxyListDownloadProxyList := controller("https://www.proxy-list.download/HTTP")
	proxyListDownloadHttps := controller("https://www.proxy-list.download/HTTPS")
	proxyListDownloadSocks4 := controller("https://www.proxy-list.download/SOCKS4")
	proxyListDownloadSocks5 := controller("https://www.proxy-list.download/SOCKS5")
	anonymousProxyOrg := controller("https://free-proxy-list.net/anonymous-proxy.html")
	ukProxy := controller("https://free-proxy-list.net/uk-proxy.html")
	usProxy := controller("https://www.us-proxy.org/")

	proxies := make([]Proxy, 0)
	proxies = append(proxies, socksProxyNet...)
	proxies = append(proxies, sslProxiesOrg...)
	proxies = append(proxies, freeProxyListNet...)
	proxies = append(proxies, usProxyOrg...)
	proxies = append(proxies, proxyNovaCom...)
	proxies = append(proxies, proxyListOrg...)
	proxies = append(proxies, proxyListDownloadProxyList...)
	proxies = append(proxies, proxyListDownloadHttps...)
	proxies = append(proxies, proxyListDownloadSocks4...)
	proxies = append(proxies, proxyListDownloadSocks5...)
	proxies = append(proxies, anonymousProxyOrg...)
	proxies = append(proxies, ukProxy...)
	proxies = append(proxies, usProxy...)

	return proxies
}

func freeProxyListCC() []Proxy {
	body, err := Dialer.Dial("https://freeproxylist.cc/servers/1.html")
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(body)

	proxies := make([]Proxy, 0)
	doc, err := goquery.NewDocumentFromReader(reader)

	var lastPage int

	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
		last, _ := doc.Find("ul.pagination li:last-child > a[href]").Attr("href")
		regex := regexp.MustCompile(`https://freeproxylist\.cc/servers/(\d+)\.html`)
		page := regex.FindStringSubmatch(last)[1]
		lastPage, _ = strconv.Atoi(page)
	})

	for i := 1; i <= lastPage; i++ {
		body, err := Dialer.Dial(fmt.Sprintf("https://freeproxylist.cc/servers/%d.html", i))
		if err != nil {
			panic(err)
		}

		reader := bytes.NewReader(body)

		doc, err := goquery.NewDocumentFromReader(reader)

		doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
			host := s.Find("td:first-child").Text()
			port := s.Find("td:nth-child(2)").Text()
			atoi, _ := strconv.Atoi(port)

			var protocolStr string
			if strings.Contains(s.Find("td:nth-child(6)").Text(), "no") {
				protocolStr = "http"
			} else {
				protocolStr = "https"
			}

			countryStr := s.Find("td:nth-child(5)").Text()
			anonymityStr := s.Find("td:nth-child(6)").Text()

			var protocol Protocol
			var anonymity Anonymity

			protocol.Parse(protocolStr)
			anonymity.Parse(anonymityStr)

			proxies = append(proxies, Proxy{
				Address:   host + ":" + port,
				Host:      host,
				Port:      uint8(atoi),
				Country:   countryStr,
				Anonymity: anonymity,
				Protocol:  protocol,
			})
		})
	}

	return proxies
}
