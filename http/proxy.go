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

func (p *proxy) SetProxy(transport *http.Transport) error {
	if !p.Ativo {
		return nil // Proxy desabilitado, não configura nada
	}

	if p.Host == "" {
		return fmt.Errorf("host do proxy não pode estar vazio")
	}

	if p.Port <= 0 || p.Port > 65535 {
		return fmt.Errorf("porta do proxy deve estar entre 1 e 65535")
	}

	// Create a new proxy URL with authentication
	proxyURL, err := url.Parse(p.getUrl())
	if err != nil {
		return fmt.Errorf("erro ao fazer parse da URL do proxy: %v", err)
	}

	transport.Proxy = http.ProxyURL(proxyURL)
	return nil
}

// SetProxyConfig configura os parâmetros básicos do proxy
func (p *proxy) SetProxyConfig(host string, port int, username, password string) {
	p.Host = host
	p.Port = port
	p.UserName = username
	p.Password = password
	p.Ativo = true
}

// SetAtivo ativa ou desativa o proxy
func (p *proxy) SetAtivo(ativo bool) {
	p.Ativo = ativo
}

// GetAtivo verifica se o proxy está ativo
func (p *proxy) GetAtivo() bool {
	return p.Ativo
}

// SetHost define o host do proxy
func (p *proxy) SetHost(host string) {
	p.Host = host
}

// GetHost retorna o host do proxy
func (p *proxy) GetHost() string {
	return p.Host
}

// SetPort define a porta do proxy
func (p *proxy) SetPort(port int) {
	p.Port = port
}

// GetPort retorna a porta do proxy
func (p *proxy) GetPort() int {
	return p.Port
}

// SetUserName define o usuário para autenticação
func (p *proxy) SetUserName(username string) {
	p.UserName = username
}

// GetUserName retorna o usuário configurado
func (p *proxy) GetUserName() string {
	return p.UserName
}

// SetPassword define a senha para autenticação
func (p *proxy) SetPassword(password string) {
	p.Password = password
}

// GetPassword retorna a senha configurada
func (p *proxy) GetPassword() string {
	return p.Password
}

func NewProxy() *proxy {
	px := &proxy{
		Ativo: false,
	}
	return px
}
