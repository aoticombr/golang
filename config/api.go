package config

import "strconv"

type Api struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Path string `json:"path"`

	Swagger Swagger  `json:"swagger"`
	Cors    Cors     `json:"cors"`
	Https   Https    `json:"https"`
	Gateway Gateway  `json:"gateway"`
	Dbs     []string `json:"dbs"`
	Ativo   bool     `json:"ativo"`
}

func (a *Api) GetPortStr() string {
	return strconv.Itoa(a.Port)
}
