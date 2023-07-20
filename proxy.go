package goproxyscrape

type Proxy struct {
	Address   string    `json:"address"` // Mix of host and port
	Host      string    `json:"host"`
	Port      uint8     `json:"port"`
	Country   string    `json:"country"`
	Anonymity Anonymity `json:"anonymity"`
	Protocol  Protocol  `json:"protocol"`
}
