package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Response struct {
	StatusCode    int
	StatusMessage string
	Body          []byte
	Header        map[string][]string
}

func (R *Response) GetStatusCodeStr() string {
	return fmt.Sprintf("%d", R.StatusCode)
}
func (R *Response) GetStatusMessage() string {
	return R.StatusMessage
}
func (R *Response) GetBody() []byte {
	return R.Body
}
func (R *Response) GetHeader() map[string][]string {
	return R.Header
}

func (R *Response) GetAllFields() map[string]string {
	headerValues := make(map[string]string)
	// Adicionando os campos extras do cabeçalho (ExtraFields)
	for fieldName, fieldValues := range R.Header {
		if len(fieldValues) > 0 {
			// Se houver múltiplos valores para o campo, concatenamos eles em uma única string
			headerValues[fieldName] = strings.Join(fieldValues, "; ")
		}
	}

	return headerValues
}

func (R *Response) GetToken() (string, error) {
	var (
		TokenResponse TokenResponse
	)

	//	fmt.Println("passou aqui a, 3", Resp.StatusCode)
	if R.StatusCode < 200 || R.StatusCode >= 300 {
		//	fmt.Println("passou aqui a, 3", Resp.StatusCode)
		//	fmt.Println("passou aqui b, 3", Resp.StatusMessage)
		return "", errors.New(fmt.Sprintf("Erro de validação de token OUTH2: %d %s", R.StatusCode, R.StatusMessage))
	} else {
		//fmt.Println("body:", string(Resp.Body))
		err := json.Unmarshal(R.Body, &TokenResponse)
		if err != nil {
			return "", err
		}
		//	fmt.Println("send.. auth...token 2")
		return TokenResponse.AccessToken, nil
	}
}
func NewResponse() *Response {
	return &Response{}
}
