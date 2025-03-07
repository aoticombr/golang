package jsonconfig

type Api struct {
	Name      string   `json:"name"`
	Protocolo string   `json:"protocolo"`
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	Gateway   Gateway  `json:"gateway"`
	Schemas   []Schema `json:"schemas"`
	Ativo     bool     `json:"ativo"`
}

func (a *Api) GetGateway() *Gateway {
	return &a.Gateway
}
func (a *Api) GetSchemas() []Schema {
	return a.Schemas
}
func (a *Api) GetSchemaByName(name string) *Schema {
	for _, schema := range a.Schemas {
		if schema.Schema == name {
			return &schema
		}
	}
	return nil
}
