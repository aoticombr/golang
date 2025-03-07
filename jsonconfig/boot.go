package jsonconfig

type Boot struct {
	Name    string   `json:"name"`
	Schemas []Schema `json:"schemas"`
	Ativo   bool     `json:"ativo"`
}

func (b *Boot) GetSchemas() []Schema {
	return b.Schemas
}
func (b *Boot) GetSchemaByName(name string) *Schema {
	for _, schema := range b.Schemas {
		if schema.Schema == name {
			return &schema
		}
	}
	return nil
}
func (b *Boot) GetSchemaBySID(sid string) *Schema {
	for _, schema := range b.Schemas {
		if schema.SID == sid {
			return &schema
		}
	}
	return nil
}
func (b *Boot) GetSchemaByHost(host string) *Schema {
	for _, schema := range b.Schemas {
		if schema.Host == host {
			return &schema
		}
	}
	return nil
}
func (b *Boot) GetSchemaByPort(port int) *Schema {
	for _, schema := range b.Schemas {
		if schema.Port == port {
			return &schema
		}
	}
	return nil
}
func (b *Boot) GetSchemaByUser(user string) *Schema {
	for _, schema := range b.Schemas {
		if schema.User == user {
			return &schema
		}
	}
	return nil
}
func NewBoot() *Boot {
	return &Boot{}
}
