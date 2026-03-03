package config

type Gateway struct {
	Protocolo string `json:"protocolo"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Ativo     bool   `json:"ativo"`
}
