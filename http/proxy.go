package http

import (
	"fmt"
	"net/http"
	"net/url"
)

type proxy struct {
	UserName string
	Password string
	Host     string
	Port     int
	Ativo    bool
}

func (p *proxy) getUrl() string {
	var auth string
	if p.UserName != "" && p.Password != "" {
		auth = p.UserName + ":" + p.Password + "@"
	}

	proxyURL := fmt.Sprintf("http://%s%s:%d", auth, p.Host, p.Port)
	return proxyURL
}

func (p *proxy) GetTransport() (*http.Transport, error) {
	if p.Ativo {
		// Create a new proxy URL with authentication
		proxy, err := url.Parse(p.getUrl())
		if err != nil {
			return nil, fmt.Errorf("Error parsing proxy URL:" + err.Error())
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		return transport, nil
	} else {
		return nil, nil
	}
}

func NewProxy() *proxy {
	px := &proxy{}
	return px
}
