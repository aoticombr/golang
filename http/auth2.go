package http

import (
	"encoding/json"
	"fmt"
)

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
	Owner        *THttp
	AuthUrl      string
	ClientId     string
	ClientSecret string
	Scope        string
	ClientAuth   ClientAuth
	Resp         *Response
	Erro         error
}

func (A *auth2) Send() (RES *Response, err error) {
	HttpToken := NewHttp()
	HttpToken.Request.Header.ContentType = "application/x-www-form-urlencoded"
	HttpToken.Request.Header.Accept = "*/*"

	HttpToken.SetUrl(A.AuthUrl)
	//fmt.Println("A.AuthUrl...", A.AuthUrl)
	HttpToken.Metodo = M_POST
	if A.Owner != nil {
		HttpToken.Certificate.PathCrt = A.Owner.Certificate.PathCrt
		HttpToken.Certificate.PathPriv = A.Owner.Certificate.PathPriv
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
	return HttpToken.Send()
}

func (A *auth2) GetToken() (string, error) {
	var (
		TokenResponse TokenResponse
	)

	Resp, err := A.Send()

	A.Resp = Resp
	A.Erro = err
	//fmt.Println("passou aqui a, 1", Resp)
	if err != nil {
		//	fmt.Println("passou aqui a, 2", err)
		return "", err
	}
	//	fmt.Println("passou aqui a, 3", Resp.StatusCode)
	if Resp.StatusCode < 200 || Resp.StatusCode >= 300 {
		//	fmt.Println("passou aqui a, 3", Resp.StatusCode)
		//	fmt.Println("passou aqui b, 3", Resp.StatusMessage)
		return "", fmt.Errorf("Erro de validação de token OUTH2", Resp.StatusCode, Resp.StatusMessage, err)
	} else {
		//fmt.Println("body:", string(Resp.Body))
		err = json.Unmarshal(Resp.Body, &TokenResponse)
		if err != nil {
			return "", err
		}
		//	fmt.Println("send.. auth...token 2")
		return TokenResponse.AccessToken, nil

	}
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
