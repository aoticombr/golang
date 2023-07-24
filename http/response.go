package http

import (
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
