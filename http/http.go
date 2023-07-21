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
	Auth2             auth2
	Request           *Request
	Metodo            TMethod
	AuthorizationType AuthorizationType
	Authorization     string
	Password          string
	UserName          string
	url               string
	Protocolo         string // http, https
	Host              string // www.example.com
	Path              string // /product
	Varibles          Varibles
	Params            Params
}

func (H *THttp) SetMetodoStr(value string) error {
	H.Metodo, _ = GetStrFromMethod(value)
	return nil
}
func (H *THttp) GetMetodoStr() string {
	return GetMethodStr(H.Metodo)
}
func (H *THttp) SetMetodo(value TMethod) error {
	H.Metodo = value
	return nil
}
func (H *THttp) GetMetodo() TMethod {
	return H.Metodo
}
func (H *THttp) SetAuthorizationType(value AuthorizationType) error {
	H.AuthorizationType = value
	return nil
}
func (H *THttp) GetAuthorizationType() AuthorizationType {
	return H.AuthorizationType
}

func (H *THttp) SetUrl(value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return err
	}
	for key, values := range u.Query() {
		H.Params.Add(key, strings.Join(values, ", "))
	}
	H.Protocolo = u.Scheme
	H.Host = u.Host
	H.Path = u.Path
	H.url = fmt.Sprintf("%s://%s%s", H.Protocolo, H.Host, H.Path)
	return nil
}
func (H *THttp) GetFullURL() (string, error) {
	return fmt.Sprintf("%s://%s", H.Protocolo, H.Host), nil
}
func (H *THttp) GetUrl() string {

	queryParams := url.Values{}

	baseURL := H.url
	//fmt.Println("baseURL:", baseURL)
	for key, value := range H.Params {
		queryParams.Add(key, value)
	}
	if strings.Contains(baseURL, "?") {
		baseURL += "&" + queryParams.Encode()
	} else {
		baseURL += "?" + queryParams.Encode()
	}
	for key, value := range H.Varibles {
		baseURL = strings.ReplaceAll(baseURL, "{{"+key+"}}", value)
	}
	return baseURL
}
func (H *THttp) CompletHeader() {
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
func (H *THttp) CompletAutorization() error {
	if H.AuthorizationType == AT_AutoDetect {
		if H.Authorization != "" {
			H.AuthorizationType = AT_Bearer
		} else if H.UserName != "" && H.Password != "" {
			H.AuthorizationType = AT_Basic
		}
	}
	if H.AuthorizationType == AT_Auth2 {
		token, err := H.Auth2.GetToken()
		if err != nil {
			return fmt.Errorf("Erro ao obter o token:", err.Error())
		}
		H.Authorization = token
	}
	if H.AuthorizationType == AT_Bearer {
		H.req.Header.Set("Authorization", "Bearer "+H.Authorization)
	}
	if H.AuthorizationType == AT_Basic {
		H.req.SetBasicAuth(H.UserName, H.Password)
	}
	return nil
}
func (H *THttp) Send() (*Response, error) {
	var err error
	var resp *http.Response
	client := &http.Client{}
	uri := H.GetUrl()
	if strings.Contains(uri, "{{") || strings.Contains(uri, "}}") {
		return nil, fmt.Errorf("Erro ao validar url, variaveis não substituidas:", uri, err)
	}
	switch GetContentTypeFromString(H.Request.Header.ContentType) {
	case CT_NONE:
		//	fmt.Println("CT_NONE:")
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), nil)
	case CT_TEXT, CT_JAVASCRIPT, CT_JSON, CT_HTML, CT_XML:
		//	fmt.Println("CT_TEXT:")
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), bytes.NewReader(H.Request.Body))
	case CT_MULTIPART_FORM_DATA:
		//fmt.Println("CT_MULTIPART_FORM_DATA:")
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
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), &requestBody)
		// Defina o cabeçalho da requisição para indicar que está enviando dados com o formato multipart/form-data
		H.Request.Header.ContentType = multipartWriter.FormDataContentType()
	case CT_X_WWW_FORM_URLENCODED:
		//	fmt.Println("CT_X_WWW_FORM_URLENCODED:")
		formData := url.Values{}
		if H.Request.ItensFormField != nil {
			for _, v := range H.Request.ItensFormField {
				formData.Add(v.FieldName, v.FieldValue)
			}
		}
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), strings.NewReader(formData.Encode()))
	case CT_BINARY:
		//fmt.Println("CT_BINARY:")
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
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), fileBuffer)
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

	ht := &THttp{
		Request:           NewRequest(),
		Params:            NewParams(),
		Varibles:          NewVaribles(),
		Metodo:            M_GET,
		AuthorizationType: AT_AutoDetect,
	}
	return ht
}
