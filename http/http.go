package http

import (
	"bytes"
	"fmt"
	"net/http"
)

type THttp struct {
	req *http.Request

	Request           *Request
	Response          []byte
	ContentType       TContentType
	Metodo            TMethod
	Charset           string
	AuthorizationType AuthorizationType
	Authorization     string
	Password          string
	UserName          string
	Url               string
}

func (H THttp) CompletHeader() {
	if H.Request.Header.Accept != "" {
		H.req.Header.Set("Accept", H.Request.Header.Accept)
	}
	if H.Request.Header.AcceptCharset != "" {
		H.req.Header.Set("Accept-Charset", H.Request.Header.AcceptCharset)
	}
	if H.Request.Header.AcceptEncoding != "" {
		H.req.Header.Set("Accept-Encoding", H.Request.Header.AcceptEncoding)
	}
	if H.Request.Header.AcceptLanguage != "" {
		H.req.Header.Set("Accept-Language", H.Request.Header.AcceptLanguage)
	}
	if H.Request.Header.Authorization != "" {
		H.req.Header.Set("Authorization", H.Request.Header.Authorization)
	}
	if H.Request.Header.Charset != "" {
		H.req.Header.Set("Charset", H.Request.Header.Charset)
	}
	if H.Request.Header.ContentType != "" {
		H.req.Header.Set("Content-Type", H.Request.Header.ContentType)
	}
	if H.Request.Header.ContentLength != "" {
		H.req.Header.Set("Content-Length", H.Request.Header.ContentLength)
	}
	if H.Request.Header.ContentEncoding != "" {
		H.req.Header.Set("Content-Encoding", H.Request.Header.ContentEncoding)
	}
	if H.Request.Header.ContentVersion != "" {
		H.req.Header.Set("Content-Version", H.Request.Header.ContentVersion)
	}
	if H.Request.Header.ContentLocation != "" {
		H.req.Header.Set("Content-Location", H.Request.Header.ContentLocation)
	}
	if H.Request.Header.ItensFormField != nil {
		for _, v := range H.Request.Header.ItensFormField {
			H.req.Form.Add(v.FieldName, v.FieldValue)
		}
	}

}
func (H THttp) CompletAutorization() {
	if H.AuthorizationType == AutoDetect {
		if H.Authorization != "" {
			H.AuthorizationType = Bearer
		} else if H.UserName != "" && H.Password != "" {
			H.AuthorizationType = Basic
		}
	}
	if H.AuthorizationType == Bearer {
		H.req.Header.Set("Authorization", "Bearer "+H.Authorization)
	}
	if H.AuthorizationType == Basic {
		H.req.SetBasicAuth(H.UserName, H.Password)
	}
}

func (H THttp) Send() error {
	var err error
	var resp *http.Response
	client := &http.Client{}
	H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.Url, bytes.NewReader(H.Response))
	if err != nil {
		return fmt.Errorf("Erro ao criar a requisição %s: %s\n", GetMethodStr(H.Metodo), err)
	}
	H.CompletHeader()
	H.CompletAutorization()

	resp, err = client.Do(H.req)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Erro ao fazer a requisição %s: %s\n", GetMethodStr(H.Metodo), err)
	}

	return nil
}

func NewHttp() *THttp {
	ht := &THttp{
		Request: &Request{
			Header: &Header{
				Accept: "*/*",
			},
		},
		Metodo:            M_GET,
		AuthorizationType: AutoDetect,
	}
	return ht
}
