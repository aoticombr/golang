package config

type Cors struct {
	MaxAge           int      `json:"maxAge"`
	AllowCredentials bool     `json:"allowCredentials"`
	AllowHeaders     []string `json:"allowHeaders"`
	ExposedHeaders   []string `json:"exposedHeaders"`
	AllowMethods     []string `json:"allowMethods"`
	AllowOrigins     []string `json:"allowOrigins"`
	Ativo            bool     `json:"ativo"`
}
