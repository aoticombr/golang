package config

type Service struct {
	Name  string   `json:"name"`
	Dbs   []string `json:"dbs"`
	Ativo bool     `json:"ativo"`
}
