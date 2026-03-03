package config

type Jwt struct {
	Name           string `json:"name"`
	ExpirationTime int    `json:"expirationTime"`
	SecretKey      string `json:"secret"`
}
