package http

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	Scope            string `json:"scope"`
}

type ClientAuth int

const (
	CA_SendBasicAuthHeader ClientAuth = iota
	CA_SendClientCredentialsInBody
)

type auth2 struct {
	Owner                   *THttp
	AuthUrl                 string
	ClientId                string
	ClientSecret            string
	Scope                   string
	ClientAuth              ClientAuth
	UseSameCertificateOwner bool
	InsecureSkipVerify      bool
}

func (A *auth2) Send() (RES *Response, err error) {
	HttpToken := NewHttp()
	HttpToken.Request.Header.ContentType = "application/x-www-form-urlencoded"
	HttpToken.Request.Header.Accept = "*/*"
	HttpToken.InsecureSkipVerify = A.InsecureSkipVerify
	HttpToken.SetUrl(A.AuthUrl)
	//fmt.Println("A.AuthUrl...", A.AuthUrl)
	HttpToken.Metodo = M_POST
	if A.Owner != nil {
		if A.UseSameCertificateOwner {
			HttpToken.Certificate = A.Owner.Certificate
		}
	}
	/*
		No OAuth2, existem 2 formas de enviar as credenciais do cliente (client_id e client_secret):

		1. Enviar as credenciais no cabeçalho de autorização usando o esquema de autenticação básica (Basic Authentication).
		   Nesse caso, o client_id e o client_secret são codificados em Base64 e enviados no cabeçalho Authorization.
		2. Enviar as credenciais no corpo da requisição usando o tipo de conteúdo application/x-www-form-urlencoded.
		   Nesse caso, o client_id e o client_secret são incluídos como parâmetros no corpo da requisição.

		A escolha entre essas duas formas depende das preferências do servidor de autorização e das práticas recomendadas de segurança.
		O envio das credenciais no cabeçalho de autorização é geralmente considerado mais seguro,
		pois as credenciais não ficam expostas no corpo da requisição. No entanto, alguns servidores de autorização
		podem exigir que as credenciais sejam enviadas no corpo da requisição, especialmente se o servidor
		não suportar autenticação básica.
	*/
	if A.ClientAuth == CA_SendBasicAuthHeader {
		HttpToken.AuthorizationType = AT_Basic
		HttpToken.UserName = A.ClientId
		HttpToken.Password = A.ClientSecret
		HttpToken.Request.AddFormField("grant_type", "client_credentials")
		if A.Scope != "" {
			HttpToken.Request.AddFormField("scope", A.Scope)
		}
	} else {
		HttpToken.AuthorizationType = AT_Nenhum
		HttpToken.Request.AddFormField("grant_type", "client_credentials")
		if A.ClientId != "" {
			HttpToken.Request.AddFormField("client_id", A.ClientId)
		}
		if A.ClientSecret != "" {
			HttpToken.Request.AddFormField("client_secret", A.ClientSecret)
		}
		if A.Scope != "" {
			HttpToken.Request.AddFormField("scope", A.Scope)
		}
	}
	HttpToken.EncType = ET_X_WWW_FORM_URLENCODED
	return HttpToken.Send()
}

func (A *auth2) GetToken() (string, error) {
	Resp, err := A.Send()
	if err != nil {
		return "", err
	}

	return Resp.GetToken()
}

func NewAuth2() *auth2 {
	return &auth2{
		AuthUrl:      "",
		ClientId:     "",
		ClientSecret: "",
		Scope:        "",
		ClientAuth:   CA_SendBasicAuthHeader,
	}
}
