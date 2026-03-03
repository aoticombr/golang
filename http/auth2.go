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
	Resp                    *Response
	Erro                    error
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
	//fmt.Println("send.. auth...token 1")
	A.Resp, A.Erro = HttpToken.Send()
	return A.Resp, A.Erro
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
