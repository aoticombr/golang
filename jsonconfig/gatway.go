package jsonconfig

type Gateway struct {
	Protocolo string `json:"protocolo"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Ativo     bool   `json:"ativo"`
}

func NewGateway() *Gateway {
	return &Gateway{}
}
