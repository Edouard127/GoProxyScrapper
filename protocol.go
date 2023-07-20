package goproxyscrape

import "fmt"

type Protocol uint8

const (
	HTTP Protocol = iota
	HTTPS
	SOCKS4
	SOCKS5
	UnknownProtocol
)

func (p Protocol) String() string {
	if p < HTTP || p > UnknownProtocol {
		return "Protocol(" + string(p) + ")"
	}
	return [...]string{"HTTP", "HTTPS", "SOCKS4", "SOCKS5", "Unknown"}[p]
}

func (p *Protocol) Parse(str string) error {
	switch str {
	case "HTTP":
		*p = HTTP
	case "HTTPS":
		*p = HTTPS
	case "SOCKS4":
		*p = SOCKS4
	case "SOCKS5":
		*p = SOCKS5
	case "Unknown":
		*p = UnknownProtocol
	default:
		return fmt.Errorf("invalid protocol: %s", str)
	}
	return nil
}
