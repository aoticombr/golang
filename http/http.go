package http

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type THttp struct {
	req               *http.Request
	Request           *Request
	Metodo            TMethod
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

	if H.Request.Header.ExtraFields != nil {
		for k, v := range H.Request.Header.ExtraFields {
			for _, v2 := range v {
				H.req.Header.Add(k, v2)
			}
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

func (H THttp) Send() (*Response, error) {
	var err error
	var resp *http.Response
	client := &http.Client{}

	switch GetContentTypeFromString(H.Request.Header.ContentType) {
	case CT_NONE:
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.Url, nil)
	case CT_TEXT, CT_JAVASCRIPT, CT_JSON, CT_HTML, CT_XML:
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.Url, bytes.NewReader(H.Request.Body))
	case CT_MULTIPART_FORM_DATA:
		var requestBody bytes.Buffer
		multipartWriter := multipart.NewWriter(&requestBody)
		defer multipartWriter.Close()
		if H.Request.ItensContentText != nil {
			for _, v := range H.Request.ItensContentText {
				multipartWriter.WriteField(v.Name, v.Value.Text())
			}
		}
		if H.Request.ItensContentBin != nil {
			for _, v := range H.Request.ItensContentBin {
				fileWriter, err := multipartWriter.CreateFormFile(v.Name, v.FileName)
				if err != nil {
					return nil, fmt.Errorf("Erro ao criar o arquivo %s: %s\n", v.FileName, err)
				}
				_, err = fileWriter.Write(v.Value)
				if err != nil {
					return nil, fmt.Errorf("Erro ao escrever o arquivo %s: %s\n", v.FileName, err)
				}
			}
		}
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.Url, &requestBody)
		// Defina o cabeçalho da requisição para indicar que está enviando dados com o formato multipart/form-data
		H.Request.Header.ContentType = multipartWriter.FormDataContentType()
	case CT_X_WWW_FORM_URLENCODED:
		formData := url.Values{}
		if H.Request.ItensFormField != nil {
			for _, v := range H.Request.ItensFormField {
				formData.Add(v.FieldName, v.FieldValue)
			}
		}
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.Url, strings.NewReader(formData.Encode()))
	case CT_BINARY:
		fileBuffer := &bytes.Buffer{}
		fileBuffer.Reset()
		if H.Request.ItensContentBin != nil {
			for _, v := range H.Request.ItensContentBin {
				_, err := fileBuffer.Write(v.Value)
				if err != nil {
					fmt.Println("Erro ao copiar os dados para o buffer:", err)
					return nil, fmt.Errorf("Erro ao copiar os dados para o buffer:", v.FileName, err)
				}
			}
		}
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.Url, fileBuffer)
	}

	if err != nil {
		return nil, fmt.Errorf("Erro ao criar a requisição %s: %s\n", GetMethodStr(H.Metodo), err)
	}
	H.CompletHeader()
	H.CompletAutorization()
	resp, err = client.Do(H.req)

	if err != nil {
		return nil, fmt.Errorf("Erro ao fazer a requisição %s: %s\n", GetMethodStr(H.Metodo), err)
	}
	defer resp.Body.Close()
	// Ler a resposta (opcional)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler body : %s\n", err)
	}
	RES := &Response{
		StatusCode:    resp.StatusCode,
		StatusMessage: resp.Status,
		Body:          body,
		Header:        resp.Header,
	}
	return RES, nil
}

func NewHttp() *THttp {
	fmt.Println("NewHttp")
	ht := &THttp{
		Request:           NewRequest(),
		Metodo:            M_GET,
		AuthorizationType: AutoDetect,
	}
	return ht
}
