package config

import "time"

type Https struct {
	Cert     string    `json:"cert"`
	Key      string    `json:"key"`
	Validate time.Time `json:"validate"`
	Ativo    bool      `json:"ativo"`
}
