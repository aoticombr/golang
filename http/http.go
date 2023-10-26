package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type THttp struct {
	req               *http.Request
	ws                *websocket.Conn
	Auth2             *auth2
	Request           *Request
	Response          *Response
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
	Proxy             *proxy
	EncType           EncType
	Timeout           int //segundos
	OnSend            IWebsocket
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
	//	fmt.Println("baseURL:", baseURL)
	if queryParams.Encode() != "" {
		if strings.Contains(baseURL, "?") {
			baseURL += "&" + queryParams.Encode()
		} else {
			baseURL += "?" + queryParams.Encode()
		}
	}
	//fmt.Println("baseURL:", baseURL)
	for key, value := range H.Varibles {
		baseURL = strings.ReplaceAll(baseURL, "{{"+key+"}}", value)
	}
	return baseURL
}
func (H *THttp) completHeader() {
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
func (H *THttp) completAutorization(req *http.Request) error {
	//	fmt.Println("passou aqui 1")
	if H.AuthorizationType == AT_AutoDetect {
		//	fmt.Println("passou aqui 1.1")
		if H.Authorization != "" {
			H.AuthorizationType = AT_Bearer
		} else if H.UserName != "" && H.Password != "" {
			H.AuthorizationType = AT_Basic
		}
	}
	//fmt.Println("passou aqui 2:>", H.AuthorizationType)
	if H.AuthorizationType == AT_Auth2 {
		//	fmt.Println("passou aqui 2.1")
		token, err := H.Auth2.GetToken()
		if err != nil {
			//fmt.Println("Erro ao obter o token:", err.Error())
			return fmt.Errorf("Erro ao obter o token:", err.Error())
		}
		H.AuthorizationType = AT_Bearer
		H.Authorization = token
		//	fmt.Println("passou aqui 3.a", "H.Authorization "+H.Authorization)
	}
	//	fmt.Println("passou aqui 3")
	if H.AuthorizationType == AT_Bearer {
		//fmt.Println("passou aqui 3.1", "Bearer "+H.Authorization)
		inputStringLower := strings.ToLower(H.Authorization)
		searchTermLower := "bearer"

		if strings.Contains(inputStringLower, searchTermLower) {
			req.Header.Set("Authorization", H.Authorization)
		} else {
			req.Header.Set("Authorization", "Bearer "+H.Authorization)
		}

	}
	//fmt.Println("passou aqui 4")
	if H.AuthorizationType == AT_Basic {
		//fmt.Println("passou aqui 5", H.UserName, H.Password)
		auth := H.UserName + ":" + H.Password
		basic := base64.StdEncoding.EncodeToString([]byte(auth))
		H.Request.Header.Authorization = "Basic " + basic

		//fmt.Println("H.Request.Header.Authorization:", H.Request.Header.Authorization)
		req.SetBasicAuth(H.UserName, H.Password)
	}
	return nil
}
func (H *THttp) completAutorizationSocket(req http.Header) error {
	//	fmt.Println("passou aqui 1")
	if H.AuthorizationType == AT_AutoDetect {
		//	fmt.Println("passou aqui 1.1")
		if H.Authorization != "" {
			H.AuthorizationType = AT_Bearer
		} else if H.UserName != "" && H.Password != "" {
			H.AuthorizationType = AT_Basic
		}
	}
	//	fmt.Println("passou aqui 2:>", H.AuthorizationType)
	if H.AuthorizationType == AT_Auth2 {
		fmt.Println("passou aqui 2.1")
		token, err := H.Auth2.GetToken()
		if err != nil {
			//fmt.Println("Erro ao obter o token:", err.Error())
			return fmt.Errorf("Erro ao obter o token:", err.Error())
		}
		H.AuthorizationType = AT_Bearer
		H.Authorization = token
		//fmt.Println("passou aqui 3.a", "H.Authorization "+H.Authorization)
	}
	//fmt.Println("passou aqui 3")
	if H.AuthorizationType == AT_Bearer {
		//fmt.Println("passou aqui 3.1", "Bearer "+H.Authorization)
		inputStringLower := strings.ToLower(H.Authorization)
		searchTermLower := "bearer"

		if strings.Contains(inputStringLower, searchTermLower) {
			req.Set("Authorization", H.Authorization)
		} else {
			req.Set("Authorization", "Bearer "+H.Authorization)
		}

	}
	//fmt.Println("passou aqui 4")
	if H.AuthorizationType == AT_Basic {
		//fmt.Println("passou aqui 5", H.UserName, H.Password)
		auth := H.UserName + ":" + H.Password
		basic := base64.StdEncoding.EncodeToString([]byte(auth))
		H.Request.Header.Authorization = "Basic " + basic

	}
	return nil
}
func (H *THttp) Send() (*Response, error) {
	//fmt.Println("==================")
	//fmt.Println("Send..")
	//fmt.Println("------------------")
	H.Response = NewResponse()
	var err error
	var resp *http.Response
	var trans *http.Transport
	trans, _ = H.Proxy.GetTransport()
	var client *http.Client
	if trans != nil {
		client = &http.Client{
			Transport: trans,
			Timeout:   time.Duration(H.Timeout) * time.Second,
		}
	} else {
		client = &http.Client{
			Timeout: time.Duration(H.Timeout) * time.Second,
		}
	}
	uri := H.GetUrl()
	if strings.Contains(uri, "{{") || strings.Contains(uri, "}}") {
		return nil, fmt.Errorf("Erro ao validar url, variaveis não substituidas:", uri, err)
	}
	switch H.EncType {
	case ET_NONE:
		//fmt.Println("CT_NONE:")
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), nil)

	case ET_FORM_DATA:
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
	case ET_X_WWW_FORM_URLENCODED:
		//fmt.Println("CT_X_WWW_FORM_URLENCODED:")
		formData := url.Values{}
		if H.Request.ItensFormField != nil {
			for _, v := range H.Request.ItensFormField {
				formData.Add(v.FieldName, v.FieldValue)
			}
		}
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), strings.NewReader(formData.Encode()))
	case ET_RAW:
		//fmt.Println("CT_TEXT:")
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), bytes.NewReader(H.Request.Body))
	case ET_BINARY:
		fmt.Println("CT_BINARY:")
		fileBuffer := &bytes.Buffer{}
		fileBuffer.Reset()
		if H.Request.ItensContentBin != nil {
			for _, v := range H.Request.ItensContentBin {
				_, err := fileBuffer.Write(v.Value)
				if err != nil {
					//fmt.Println("Erro ao copiar os dados para o buffer:", err)
					return nil, fmt.Errorf("Erro ao copiar os dados para o buffer:", v.FileName, err)
				}
			}
		}
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), H.GetUrl(), fileBuffer)
	}

	if err != nil {
		return nil, fmt.Errorf("Erro ao criar a requisição %s: %s\n", GetMethodStr(H.Metodo), err)
	}

	H.completAutorization(H.req)
	H.completHeader()
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
	H.Response = RES
	return RES, nil
}
func (H *THttp) Conectar() error {
	var (
		err     error
		headers http.Header
	)

	headers = make(http.Header)
	H.completAutorizationSocket(headers)
	if H.Request.Header.Accept != "" {
		headers.Set("Accept", H.Request.Header.Accept)
	}
	if H.Request.Header.AcceptCharset != "" {
		headers.Set("Accept-Charset", H.Request.Header.AcceptCharset)
	}
	if H.Request.Header.AcceptEncoding != "" {
		headers.Set("Accept-Encoding", H.Request.Header.AcceptEncoding)
	}
	if H.Request.Header.AcceptLanguage != "" {
		headers.Set("Accept-Language", H.Request.Header.AcceptLanguage)
	}
	if H.Request.Header.Authorization != "" {
		headers.Set("Authorization", H.Request.Header.Authorization)
	}
	if H.Request.Header.Charset != "" {
		headers.Set("Charset", H.Request.Header.Charset)
	}
	if H.Request.Header.ContentType != "" {
		headers.Set("Content-Type", H.Request.Header.ContentType)
	}
	if H.Request.Header.ContentLength != "" {
		headers.Set("Content-Length", H.Request.Header.ContentLength)
	}
	if H.Request.Header.ContentEncoding != "" {
		headers.Set("Content-Encoding", H.Request.Header.ContentEncoding)
	}
	if H.Request.Header.ContentVersion != "" {
		headers.Set("Content-Version", H.Request.Header.ContentVersion)
	}
	if H.Request.Header.ContentLocation != "" {
		headers.Set("Content-Location", H.Request.Header.ContentLocation)
	}

	if H.Request.Header.ExtraFields != nil {
		for k, v := range H.Request.Header.ExtraFields {
			for _, v2 := range v {
				headers.Add(k, v2)
			}
		}
	}

	H.ws, _, err = websocket.DefaultDialer.Dial(H.GetUrl(), headers)
	if err != nil {
		return err
	}
	go func() {
		for {
			messageType, p, err := H.ws.ReadMessage()
			if err != nil {
				if H.OnSend != nil {
					go H.OnSend.Read(messageType, p, err)
				}
				return
			}

		}
	}()
	return nil
}
func (H *THttp) IsConect() bool {
	if H.ws != nil {
		return true
	}
	return false
}
func (H *THttp) Desconectar() error {
	return H.ws.Close()
}
func (H *THttp) EnviarBinario(messageType int, data []byte) error {
	return H.ws.WriteMessage(messageType, data)
}
func (H *THttp) EnviarTexto(messageType int, data string) error {
	return H.ws.WriteMessage(messageType, []byte(data))
}
func (H *THttp) EnviarTextTypeTextMessage(data []byte) error {
	return H.ws.WriteMessage(websocket.TextMessage, data)
}
func (H *THttp) EnviarBinarioTypeBinaryMessage(data []byte) error {
	return H.ws.WriteMessage(websocket.BinaryMessage, data)
}

func NewHttp() *THttp {

	ht := &THttp{
		Request:           NewRequest(),
		Response:          NewResponse(),
		Params:            NewParams(),
		Varibles:          NewVaribles(),
		Proxy:             NewProxy(),
		Auth2:             NewAuth2(),
		Metodo:            M_GET,
		Timeout:           30,
		AuthorizationType: AT_AutoDetect,
	}
	return ht
}
