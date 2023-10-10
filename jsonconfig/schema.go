package jsonconfig

type Schema struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Schema string `json:"schema"`
	SID    string `json:"sid"`
	Ativo  bool   `json:"ativo"`
}

func NewSchema() *Schema {
	return &Schema{}
}
