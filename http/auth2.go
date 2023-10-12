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
	AuthUrl      string
	ClientId     string
	ClientSecret string
	Scope        string
	ClientAuth   ClientAuth
	Resp         *Response
	Erro         error
}

func (A *auth2) GetToken() (string, error) {
	var (
		TokenResponse TokenResponse
	)
	HttpToken := NewHttp()
	HttpToken.Request.Header.ContentType = "application/x-www-form-urlencoded"
	HttpToken.Request.Header.Accept = "*/*"

	HttpToken.SetUrl(A.AuthUrl)
	fmt.Println("A.AuthUrl...", A.AuthUrl)
	HttpToken.Metodo = M_POST

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
		HttpToken.Request.AddFormField("client_id", A.ClientId)
		HttpToken.Request.AddFormField("client_secret", A.ClientSecret)
		if A.Scope != "" {
			HttpToken.Request.AddFormField("scope", A.Scope)
		}
	}
	fmt.Println("send.. auth...token 1")
	Resp, err := HttpToken.Send()

	A.Resp = Resp
	A.Erro = err
	fmt.Println("passou aqui a, 1", Resp)
	if err != nil {
		fmt.Println("passou aqui a, 2", err)
		return "", err
	}
	//	fmt.Println("passou aqui a, 3", Resp.StatusCode)
	if Resp.StatusCode < 200 || Resp.StatusCode >= 300 {
		fmt.Println("passou aqui a, 3", Resp.StatusCode)
		fmt.Println("passou aqui b, 3", Resp.StatusMessage)
		return "", fmt.Errorf("Erro de validação de token OUTH2", Resp.StatusCode, Resp.StatusMessage, err)
	} else {
		fmt.Println("body:", string(Resp.Body))
		err = json.Unmarshal(Resp.Body, &TokenResponse)
		if err != nil {
			return "", err
		}
		fmt.Println("send.. auth...token 2")
		return TokenResponse.AccessToken, nil

	}

	return "", nil
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
