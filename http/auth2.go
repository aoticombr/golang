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

type auth2 struct {
	AuthUrl      string
	ClientId     string
	ClientSecret string
	Scope        string
	Resp         *Response
	Erro         error
}

func (A *auth2) GetToken() (string, error) {
	var (
		TokenResponse TokenResponse
	)
	HttpCliente := NewHttp()
	HttpCliente.Request.Header.ContentType = "application/x-www-form-urlencoded"
	HttpCliente.Request.Header.Accept = "*/*"
	HttpCliente.AuthorizationType = AT_Basic
	HttpCliente.UserName = A.ClientId
	HttpCliente.Password = A.ClientSecret
	HttpCliente.SetUrl(A.AuthUrl)
	HttpCliente.Metodo = M_POST
	HttpCliente.Request.AddFormField("grant_type", "client_credentials")
	if A.Scope != "" {
		HttpCliente.Request.AddFormField("scope", A.Scope)
	}

	Resp, err := HttpCliente.Send()
	A.Resp = Resp
	A.Erro = err
	fmt.Println("passou aqui a, 1")
	if err != nil {
		fmt.Println("passou aqui a, 2", err)
		return "", err
	}
	//	fmt.Println("passou aqui a, 3", Resp.StatusCode)
	if Resp.StatusCode < 200 || Resp.StatusCode >= 300 {
		return "", fmt.Errorf("Erro de validação de token OUTH2", Resp.StatusCode, Resp.StatusMessage, err)
	} else {
		//fmt.Println("body:", string(Resp.Body))
		err = json.Unmarshal(Resp.Body, &TokenResponse)
		if err != nil {
			return "", err
		}
		return TokenResponse.AccessToken, nil

	}
	return "", nil
}
