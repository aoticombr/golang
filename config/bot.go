package config

type Bot struct {
	Name  string   `json:"name"`
	Dbs   []string `json:"dbs"`
	Ativo bool     `json:"ativo"`
}
