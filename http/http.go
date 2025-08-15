package http

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MSG_DISCONECT    = "Perca de Conexão..."
	MSG_RECONECTANDO = "Reconectando..."
	MSG_RECONECTADO  = "Reconectado..."
	MSG_CONECTADO    = "Conectado..."
)

type THttp struct {
	/*privado*/
	req      *http.Request
	ws       *websocket.Conn
	url      string
	urlFinal string

	/*publico*/
	Auth2              *auth2
	Request            *Request
	Response           *Response
	Metodo             TMethod
	AuthorizationType  AuthorizationType
	WebSocket          *WebSocket
	Authorization      string
	Password           string
	UserName           string
	Certificate        TCert
	TransportType      TTransport
	InsecureSkipVerify bool // usado para TLS, se for true, ignora a verificação do certificado

	Protocolo string // http, https
	Host      string // www.example.com
	Path      string // /product
	Varibles  Varibles
	Params    Params
	Proxy     *proxy
	EncType   EncType
	Timeout   int //segundos
	OnSend    IWebsocket
}

func NewHttp() *THttp {

	ht := &THttp{
		Request:           NewRequest(),
		Response:          NewResponse(),
		Params:            NewParams(),
		Varibles:          NewVaribles(),
		Proxy:             NewProxy(),
		Auth2:             NewAuth2(),
		WebSocket:         NewWebSocket(),
		Metodo:            M_GET,
		Timeout:           30,
		AuthorizationType: AT_AutoDetect,
		TransportType:     TNenhum,
	}
	ht.Auth2.Owner = ht
	return ht
}

func (H *THttp) Free() {
	if H.ws != nil {
		H.ws.Close()
		H.ws = nil
	}

	H.Request = nil
	H.Response = nil
	H.Auth2 = nil
	H.WebSocket = nil
	H.Params = nil
	H.Varibles = nil
	H.Proxy = nil
	H.req = nil
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

	//fmt.Println("passou aqui 2:>", H.AuthorizationType)
	if H.AuthorizationType == AT_Auth2 {
		//	fmt.Println("passou aqui 2.1")
		token, err := H.Auth2.GetToken()
		if err != nil {
			//fmt.Println("Erro ao obter o token:", err.Error())
			return fmt.Errorf("erro ao obter o token: %v", err)
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
		//	fmt.Println("passou aqui 2.1")
		token, err := H.Auth2.GetToken()
		if err != nil {
			//fmt.Println("Erro ao obter o token:", err.Error())
			return fmt.Errorf("erro ao obter o token: %v", err)
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
		var auth string
		if H.UserName != "" && H.Password != "" {
			auth = H.UserName + ":" + H.Password
		} else if H.UserName != "" {
			auth = H.UserName
		}

		basic := base64.StdEncoding.EncodeToString([]byte(auth))
		H.Request.Header.Authorization = "Basic " + basic

	}
	return nil
}

func (H *THttp) GetUrlFinal() string {
	return H.urlFinal
}

func (H *THttp) GetTransport() *http.Transport {
	var needTransport bool
	var transport *http.Transport

	// Verificar se precisamos de um transport customizado
	if H.Proxy != nil && H.Proxy.Ativo {
		needTransport = true
	}
	if H.Certificate.PathCrt != "" && H.Certificate.PathPriv != "" {
		needTransport = true
	}
	if H.TransportType != TNenhum {
		needTransport = true
	}
	// Adicionar verificação para HTTPS como critério (removido porque será tratado no Send())
	// if strings.EqualFold(H.Protocolo, "HTTPS") {
	//     needTransport = true
	// }

	// Só criar transport se realmente precisar
	if needTransport {
		transport = &http.Transport{}

		// Configurar proxy se ativo
		if H.Proxy != nil && H.Proxy.Ativo {
			err := H.Proxy.SetProxy(transport)
			if err != nil {
				// Log do erro, mas continua sem proxy
				fmt.Printf("Erro ao configurar proxy: %v\n", err)
			}
		}

		return transport
	}

	return nil
}

func (H *THttp) Send() (RES *Response, err error) {
	defer func() {
		if r := recover(); r != nil {
			//fmt.Printf("Recuperado de um panic no método Send: %v\n", r)
			RES = nil
			err = fmt.Errorf("recuperado de um panic no método Send: %v", r)
		}
	}()
	//fmt.Println("Send..")
	//fmt.Println("------------------")
	H.Response = NewResponse()

	var resp *http.Response
	var trans *http.Transport
	var cert tls.Certificate
	var Config *tls.Config

	// Obter transport (só será criado se necessário)
	trans = H.GetTransport()

	// Configurar certificados se necessário
	if H.Certificate.PathCrt != "" && H.Certificate.PathPriv != "" {
		cert, err = tls.LoadX509KeyPair(H.Certificate.PathCrt, H.Certificate.PathPriv)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar certificado: %v", err)
		}
	}

	// Determinar se precisa de configuração TLS
	var needsTLS bool
	var isHTTPS = strings.EqualFold(H.Protocolo, "HTTPS")

	// Precisa TLS se:
	needsTLS = isHTTPS || // URL é HTTPS
		H.TransportType == TSSL || // Forçado por TransportType
		H.TransportType == TSSLTLS || // Forçado por TransportType
		(H.Certificate.PathCrt != "" && H.Certificate.PathPriv != "") // Tem certificados

	// Configurar TLS se necessário
	if needsTLS {
		Config = &tls.Config{
			InsecureSkipVerify: H.InsecureSkipVerify,
		}

		// Adicionar certificados se existirem
		if H.Certificate.PathCrt != "" && H.Certificate.PathPriv != "" {
			Config.Certificates = []tls.Certificate{cert}
		}

		// Se não temos transport mas precisamos de TLS, criar um
		if trans == nil {
			trans = &http.Transport{}
		}

		// Aplicar configuração TLS ao transport
		trans.TLSClientConfig = Config
	}

	// Criar cliente HTTP
	client := &http.Client{Timeout: time.Duration(H.Timeout) * time.Second}

	// Usar transport apenas se ele foi criado (significa que é necessário)
	if trans != nil {
		client.Transport = trans
	}

	uri := H.GetUrl()
	H.urlFinal = uri
	if strings.Contains(uri, "{{") || strings.Contains(uri, "}}") {
		return nil, fmt.Errorf("erro ao validar URL, variáveis não substituídas: %s", uri)
	}
	switch H.EncType {
	case ET_NONE:
		//fmt.Println("CT_NONE:")
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, nil)

	case ET_FORM_DATA:
		//fmt.Println("CT_MULTIPART_FORM_DATA:")
		var requestBody bytes.Buffer
		multipartWriter := NewWriter(&requestBody)
		//defer multipartWriter.Close()
		if H.Request.ItensFormField != nil {
			for _, v := range H.Request.ItensFormField {
				if v.ContentType != "" {
					fileWriter, err := multipartWriter.CreateFormFile3(v.FieldName, v.ContentType)
					if err != nil {
						return nil, fmt.Errorf("erro ao criar o arquivo %s: %s\n", v.FieldName, err)
					}
					_, err = fileWriter.Write([]byte(v.FieldValue))
					if err != nil {
						return nil, fmt.Errorf("erro ao escrever o arquivo %s: %s\n", v.FieldName, err)
					}
				} else {
					multipartWriter.WriteField(v.FieldName, v.FieldValue)
				}
			}
		}
		if H.Request.ItensContentText != nil {
			for _, v := range H.Request.ItensContentText {
				multipartWriter.WriteField(v.Name, v.Value.Text())
			}
		}
		if H.Request.ItensContentBin != nil {
			for _, v := range H.Request.ItensContentBin {
				fileWriter, err := multipartWriter.CreateFormFile(v.Name, v.FileName)
				if err != nil {
					return nil, fmt.Errorf("erro ao criar o arquivo %s: %v", v.FileName, err)
				}
				_, err = fileWriter.Write(v.Value)
				if err != nil {
					return nil, fmt.Errorf("erro ao escrever o arquivo %s: %v", v.FileName, err)
				}
			}
		}
		if H.Request.ItensSubmitFile != nil {
			for _, v := range H.Request.ItensSubmitFile {
				var (
					fileWriter io.Writer
					err        error
				)
				// boundary := multipartWriter.Boundary()
				// fileHeader := fmt.Sprintf("--%s\r\nContent-Disposition: form-data; name=\"\"; filename=\"testepaulo.pdf\"\r\nContent-Type: application/pdf\r\n\r\n", boundary)
				// requestBody.Write([]byte(fileHeader))
				// _, err = fileWriter.Write(v.Content)
				// if err != nil {
				// 	return nil, fmt.Errorf("Erro ao escrever o arquivo %s: %s\n", v.FileName, err)
				// }
				// requestBody.Write([]byte(fmt.Sprintf("\r\n--%s--\r\n", boundary)))
				if v.ContentTransferEncoding > 0 {
					fileWriter, err = multipartWriter.CreateFormFile4(v.Key, v.FileName, v.ContentType, v.ContentTransferEncoding)
				} else {
					fileWriter, err = multipartWriter.CreateFormFile2(v.Key, v.FileName, v.ContentType)
				}

				//fileWriter, err := multipartWriter.CreateFormFile(v.Key, v.FileName)
				if err != nil {
					return nil, fmt.Errorf("erro ao criar o arquivo %s: %s\n", v.FileName, err)
				}
				_, err = fileWriter.Write(v.Content)
				if err != nil {
					return nil, fmt.Errorf("erro ao escrever o arquivo %s: %s\n", v.FileName, err)
				}
			}
		}
		multipartWriter.Close() //isso aqui nao fecha e sim escreve a ultima linha
		H.Request.Header.ContentType = multipartWriter.FormDataContentType()
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, &requestBody)

		// Defina o cabeçalho da requisição para indicar que está enviando dados com o formato multipart/form-data

	case ET_X_WWW_FORM_URLENCODED:
		//fmt.Println("CT_X_WWW_FORM_URLENCODED:")
		formData := url.Values{}
		if H.Request.ItensFormField != nil {
			for _, v := range H.Request.ItensFormField {
				formData.Add(v.FieldName, v.FieldValue)
			}
		}
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, strings.NewReader(formData.Encode()))
	case ET_RAW:
		//fmt.Println("CT_TEXT:")
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, bytes.NewReader(H.Request.Body))
	case ET_BINARY:
		//fmt.Println("CT_BINARY:")
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
		H.req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, fileBuffer)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição %s: %v", GetMethodStr(H.Metodo), err)
	}
	if H.AuthorizationType == AT_AutoDetect {
		//	fmt.Println("passou aqui 1.1")
		if H.Authorization != "" {
			H.AuthorizationType = AT_Bearer
		} else if H.UserName != "" && H.Password != "" {
			H.AuthorizationType = AT_Basic
		}
	}
	if H.AuthorizationType == AT_Auth2 && H.Auth2.AuthUrl != "" && (H.Auth2.AuthUrl == H.GetUrl() || H.GetUrl() == "") {
		RES, err = H.Auth2.Send()
		if err != nil {
			return nil, fmt.Errorf("erro ao fazer a requisição %s: %v", GetMethodStr(H.Metodo), err)
		}
		H.Response = RES
		return RES, nil

	} else {
		H.completAutorization(H.req)
		H.completHeader()
		resp, err = client.Do(H.req)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao fazer a requisição %s: %v", GetMethodStr(H.Metodo), err)
	}
	defer resp.Body.Close()
	// Ler a resposta (opcional)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler body: %v", err)
	}
	RES = &Response{
		StatusCode:    resp.StatusCode,
		StatusMessage: resp.Status,
		Body:          body,
		Header:        resp.Header,
	}
	H.Response = RES
	return RES, nil
}
func (H *THttp) websocketClient() error {
	var (
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
	dialer := websocket.DefaultDialer
	var (
		//conn *websocket.Conn
		//resp *http.Response
		err  error
		err2 error
	)

	H.ws, _, err = dialer.Dial(H.GetUrl(), headers)
	if err != nil {
		if H.OnSend != nil {
			H.OnSend.Error("Erro na conexão: " + err.Error())
		} else {
			fmt.Printf("Erro na conexão: %v\n", err)
		}
		return fmt.Errorf("Erro na conexão: " + err.Error())
	} else {
		H.WebSocket.connect = OPEN
		if H.OnSend != nil {
			H.OnSend.Msg(MSG_CONECTADO)
		} else {
			fmt.Println(MSG_CONECTADO)
		}
	}
	go func() {
		for {

			//fmt.Println("################", err)
			if ((H.ws == nil) || (err2 != nil)) && (H.WebSocket.AutoReconnect == true) {
				H.WebSocket.connect = CONNECTING
				if H.WebSocket.attempts >= H.WebSocket.NumberOfAttempts {
					break
				}
				if H.OnSend != nil {
					H.OnSend.Disconect(MSG_DISCONECT, false)
					H.OnSend.Msg(MSG_RECONECTANDO)

				} else {
					fmt.Printf(MSG_RECONECTANDO)
				}
				H.ws, _, err = dialer.Dial(H.GetUrl(), headers)
				if err != nil {
					if H.OnSend != nil {
						H.OnSend.Error("Erro na conexão: " + err.Error())
					} else {
						fmt.Printf("Erro na conexão: %v\n", err)
					}
					time.Sleep(5 * time.Second)
					H.WebSocket.attempts++
					continue
				}
				H.WebSocket.connect = OPEN
				H.WebSocket.attempts = 0
				if H.OnSend != nil {
					H.OnSend.Msg(MSG_RECONECTADO)
				} else {
					fmt.Printf(MSG_RECONECTADO)
				}

			} else if ((H.ws == nil) || (err2 != nil)) && (H.WebSocket.AutoReconnect == false) {
				H.WebSocket.connect = CLOSED
				break
			}
			err2 = nil

			for err2 == nil {
				//fmt.Println("Conectado ao servidor WebSocket 2", err2)
				msgtype, msg, err := H.ws.ReadMessage()
				if err != nil {
					if H.OnSend != nil {
						H.OnSend.Error("Erro na leitura da mensagem: " + err.Error())

					} else {
						fmt.Printf("Erro na leitura da mensagem: %v\n", err)
					}
					H.ws.Close()
					time.Sleep(5 * time.Second)
					//fmt.Println("Conectado ao servidor WebSocket 4", err)
					err2 = err
					//break
					continue
				} else {
					//fmt.Println("Conectado ao servidor WebSocket 3", err2)
					if H.OnSend != nil {
						H.OnSend.Read(msgtype, msg, err)
					}
				}
				//fmt.Println("Conectado ao servidor WebSocket 5", err2)

			}
		}
		if H.OnSend != nil {
			H.OnSend.Disconect(MSG_DISCONECT, true)

		} else {
			fmt.Printf("Erro na leitura da mensagem: %v\n", err)
		}
		H.WebSocket.connect = CLOSED
	}()
	return nil
}
func (H *THttp) Conectar() error {
	//for {
	err := H.websocketClient()
	if err != nil {
		if H.OnSend != nil {
			H.OnSend.Error("Erro na conexão: " + err.Error())
			H.OnSend.Msg("Tentando reconectar em 5 segundos...")
		} else {
			fmt.Printf("Erro na conexão: " + err.Error())
			fmt.Println("Tentando reconectar em 5 segundos...")
		}
		return err
		//time.Sleep(5 * time.Second)
	}
	//}
	return nil
}
func (H *THttp) IsConect() bool {
	if H.ws != nil {
		return true
	}
	return false
}
func (H *THttp) Desconectar() error {
	if H.ws == nil {
		return fmt.Errorf("conexão WebSocket não inicializada")
	}
	return H.ws.Close()
}
func (H *THttp) EnviarBinario(messageType int, data []byte) error {
	if H.ws == nil {
		return fmt.Errorf("erro ao enviar mensagem, conexão não estabelecida")
	}
	return H.ws.WriteMessage(messageType, data)
}
func (H *THttp) EnviarTexto(messageType int, data string) error {
	if H.ws == nil {
		return fmt.Errorf("erro ao enviar mensagem, conexão não estabelecida")
	}
	return H.ws.WriteMessage(messageType, []byte(data))
}
func (H *THttp) EnviarTextTypeTextMessage(data []byte) error {
	if H.ws == nil {
		return fmt.Errorf("erro ao enviar mensagem, conexão não estabelecida")
	}
	return H.ws.WriteMessage(websocket.TextMessage, data)
}
func (H *THttp) EnviarBinarioTypeBinaryMessage(data []byte) error {
	if H.ws == nil {
		return fmt.Errorf("erro ao enviar mensagem, conexão não estabelecida")
	}
	return H.ws.WriteMessage(websocket.BinaryMessage, data)
}

func (H *THttp) ConvertBodyInStruct(value any) error {
	err := json.Unmarshal(H.Response.Body, value)
	if err != nil {
		return err
	}
	return nil
}

// Métodos para configuração de proxy

// SetProxyConfig configura o proxy básico
func (H *THttp) SetProxyConfig(host string, port int, username, password string) {
	H.Proxy.SetProxyConfig(host, port, username, password)
}

// SetProxyAtivo ativa ou desativa o uso do proxy
func (H *THttp) SetProxyAtivo(ativo bool) {
	H.Proxy.SetAtivo(ativo)
}

// GetProxyAtivo verifica se o proxy está ativo
func (H *THttp) GetProxyAtivo() bool {
	return H.Proxy.GetAtivo()
}

// SetProxyHost define o host do proxy
func (H *THttp) SetProxyHost(host string) {
	H.Proxy.SetHost(host)
}

// GetProxyHost retorna o host do proxy
func (H *THttp) GetProxyHost() string {
	return H.Proxy.GetHost()
}

// SetProxyPort define a porta do proxy
func (H *THttp) SetProxyPort(port int) {
	H.Proxy.SetPort(port)
}

// GetProxyPort retorna a porta do proxy
func (H *THttp) GetProxyPort() int {
	return H.Proxy.GetPort()
}

// SetProxyUserName define o usuário para autenticação no proxy
func (H *THttp) SetProxyUserName(username string) {
	H.Proxy.SetUserName(username)
}

// GetProxyUserName retorna o usuário do proxy
func (H *THttp) GetProxyUserName() string {
	return H.Proxy.GetUserName()
}

// SetProxyPassword define a senha para autenticação no proxy
func (H *THttp) SetProxyPassword(password string) {
	H.Proxy.SetPassword(password)
}

// GetProxyPassword retorna a senha do proxy
func (H *THttp) GetProxyPassword() string {
	return H.Proxy.GetPassword()
}
