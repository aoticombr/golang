package http

type token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type auth2 struct {
	AuthUrl      string
	ClientId     string
	ClientSecret string
	Scope        string
}

func (A *auth2) GetToken() (string, error) {
	return "", nil
}
