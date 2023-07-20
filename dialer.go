package goproxyscrape

import (
	"github.com/valyala/fasthttp"
	"sync"
)

type HTTPDialer struct {
	Client *fasthttp.Client
	mutex  *sync.Mutex
}

func NewHTTPDialer() *HTTPDialer {
	return &HTTPDialer{
		Client: &fasthttp.Client{},
		mutex:  &sync.Mutex{},
	}
}

func (d *HTTPDialer) Dial(address string) (body []byte, err error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(address)
	resp := fasthttp.AcquireResponse()
	err = d.Client.Do(req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
