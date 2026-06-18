package http

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
)

type TAwsSignature struct {
	SecretAccessKey string
	SecretKey       string
	AwsRegion       string
	ServiceName     string
	SessionToken    string
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	Scope            string `json:"scope"`
}

type GrantType int

const (
	GT_ClientCredentials GrantType = iota
	GT_PasswordCredentials
	GT_Implicit
	GT_AuthorizationCode
	GT_AuthorizationCodeWithPKCE
)

type ClientAuth int

const (
	CA_SendBasicAuthHeader ClientAuth = iota
	CA_SendClientCredentialsInBody
)

type auth2 struct {
	Owner *THttp

	AccessTokenUrl string
	AuthUrl        string
	CallBackUrl    string

	ClientId     string
	ClientSecret string
	GrantType    GrantType

	Scope      string
	ClientAuth ClientAuth

	Code string // authorization code obtido externamente (redirect do navegador)

	// PKCE (RFC 7636) - usado quando GrantType = GT_AuthorizationCodeWithPKCE.
	CodeVerifier        string // segredo aleatório; se vazio, GetAuthorizationUrl gera um
	CodeChallengeMethod string // "S256" (padrão) ou "plain"

	UseSameCertificateOwner bool
	InsecureSkipVerify      bool

	Resp *Response
	Erro error
}

func (A *auth2) Send() (RES *Response, err error) {
	HttpToken := NewHttp()
	HttpToken.Request.Header.ContentType = "application/x-www-form-urlencoded"
	HttpToken.Request.Header.Accept = "*/*"
	HttpToken.InsecureSkipVerify = A.InsecureSkipVerify
	HttpToken.SetUrl(A.AccessTokenUrl)
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
	// (a) parâmetros específicos do grant type
	switch A.GrantType {
	case GT_AuthorizationCode, GT_AuthorizationCodeWithPKCE:
		HttpToken.Request.AddFormField("grant_type", "authorization_code")
		HttpToken.Request.AddFormField("code", A.Code)
		if A.CallBackUrl != "" {
			HttpToken.Request.AddFormField("redirect_uri", A.CallBackUrl)
		}
		// PKCE: envia o code_verifier que casa com o code_challenge da autorização.
		if A.CodeVerifier != "" {
			HttpToken.Request.AddFormField("code_verifier", A.CodeVerifier)
		}
	default: // GT_ClientCredentials (e demais ainda não implementados)
		HttpToken.Request.AddFormField("grant_type", "client_credentials")
	}
	if A.Scope != "" {
		HttpToken.Request.AddFormField("scope", A.Scope)
	}

	// (b) transmissão das credenciais do cliente (ortogonal ao grant type)
	if A.ClientAuth == CA_SendBasicAuthHeader {
		HttpToken.AuthorizationType = AT_Basic
		HttpToken.UserName = A.ClientId
		HttpToken.Password = A.ClientSecret
	} else {
		HttpToken.AuthorizationType = AT_Nenhum
		if A.ClientId != "" {
			HttpToken.Request.AddFormField("client_id", A.ClientId)
		}
		if A.ClientSecret != "" {
			HttpToken.Request.AddFormField("client_secret", A.ClientSecret)
		}
	}
	HttpToken.EncType = ET_X_WWW_FORM_URLENCODED
	A.Resp, A.Erro = HttpToken.Send()
	return A.Resp, A.Erro
}

// GetAuthorizationUrl monta a URL do endpoint de autorização (AuthUrl) que o usuário
// deve abrir no navegador para autenticar e consentir. Após o login, o servidor redireciona
// para CallBackUrl com o parâmetro ?code=..., que deve ser preenchido em A.Code para então
// chamar GetToken() (grant_type=authorization_code).
//
// state é opcional (proteção CSRF / correlação); passe "" para omitir.
//
// Para GrantType = GT_AuthorizationCodeWithPKCE: se A.CodeVerifier estiver vazio, um é
// gerado e guardado em A.CodeVerifier; o code_challenge correspondente (S256 por padrão)
// é adicionado à URL. O mesmo A.CodeVerifier deve ser usado depois em GetToken().
func (A *auth2) GetAuthorizationUrl(state string) (string, error) {
	u, err := url.Parse(A.AuthUrl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", A.ClientId)
	if A.CallBackUrl != "" {
		q.Set("redirect_uri", A.CallBackUrl)
	}
	if A.Scope != "" {
		q.Set("scope", A.Scope)
	}
	if state != "" {
		q.Set("state", state)
	}
	if A.GrantType == GT_AuthorizationCodeWithPKCE {
		if A.CodeVerifier == "" {
			if _, err := A.GeneratePKCE(); err != nil {
				return "", err
			}
		}
		method := A.CodeChallengeMethod
		if method == "" {
			method = "S256"
		}
		q.Set("code_challenge", A.CodeChallenge())
		q.Set("code_challenge_method", method)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// GeneratePKCE gera um novo code_verifier aleatório (RFC 7636, 43 chars base64url) e o
// guarda em A.CodeVerifier, retornando-o. GetAuthorizationUrl chama isto automaticamente
// quando A.CodeVerifier está vazio. Se o verifier for gerado em outro processo (ex.: app
// Delphi que monta a URL), basta atribuir A.CodeVerifier com aquele mesmo valor.
func (A *auth2) GeneratePKCE() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	A.CodeVerifier = base64.RawURLEncoding.EncodeToString(b)
	return A.CodeVerifier, nil
}

// CodeChallenge devolve o code_challenge derivado de A.CodeVerifier conforme
// A.CodeChallengeMethod ("plain" usa o próprio verifier; qualquer outro valor usa S256).
// Retorna "" se A.CodeVerifier estiver vazio.
func (A *auth2) CodeChallenge() string {
	if A.CodeVerifier == "" {
		return ""
	}
	if A.CodeChallengeMethod == "plain" {
		return A.CodeVerifier
	}
	sum := sha256.Sum256([]byte(A.CodeVerifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
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
