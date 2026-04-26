package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/pkcs12"
)

const (
	MSG_DISCONECT    = "Perca de Conexão..."
	MSG_RECONECTANDO = "Reconectando..."
	MSG_RECONECTADO  = "Reconectado..."
	MSG_CONECTADO    = "Conectado..."
)

type THttp struct {
	/*privado*/
	ws       *websocket.Conn
	wsMu     sync.RWMutex // protege ws e o estado interno de WebSocket (connect, attempts) acessado pela goroutine de leitura
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

// Free libera explicitamente todos os recursos internos para o GC.
// Usado quando o componente recebeu payloads grandes e precisamos garantir
// que tudo seja recuperado de imediato (não confiamos só no escopo do caller).
func (H *THttp) Free() {
	H.wsMu.Lock()
	if H.ws != nil {
		H.ws.Close()
		H.ws = nil
	}
	H.wsMu.Unlock()

	H.Request = nil
	H.Response = nil
	H.Auth2 = nil
	H.WebSocket = nil
	H.Params = nil
	H.Varibles = nil
	H.Proxy = nil
}

func (H *THttp) SetMetodoStr(value string) error {
	var err error
	H.Metodo, err = GetStrFromMethod(value)
	return err
}
func (H *THttp) GetMetodoStr() string {
	return GetMethodStr(H.Metodo)
}
func (H *THttp) SetMetodo(value TMethod) {
	H.Metodo = value
}
func (H *THttp) GetMetodo() TMethod {
	return H.Metodo
}
func (H *THttp) SetAuthorizationType(value AuthorizationType) {
	H.AuthorizationType = value
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

// applyHeaders copia os campos não-vazios de H.Request.Header para um http.Header
// destino. Compartilhado entre completHeader() (HTTP) e websocketClient() (WS) para
// evitar 30+ linhas duplicadas.
func (H *THttp) applyHeaders(dst http.Header) {
	if H.Request == nil {
		return
	}
	h := H.Request.Header
	if h.Accept != "" {
		dst.Set("Accept", h.Accept)
	}
	if h.AcceptCharset != "" {
		dst.Set("Accept-Charset", h.AcceptCharset)
	}
	if h.AcceptEncoding != "" {
		dst.Set("Accept-Encoding", h.AcceptEncoding)
	}
	if h.AcceptLanguage != "" {
		dst.Set("Accept-Language", h.AcceptLanguage)
	}
	if h.Authorization != "" {
		dst.Set("Authorization", h.Authorization)
	}
	if h.Charset != "" {
		dst.Set("Charset", h.Charset)
	}
	if h.ContentType != "" {
		dst.Set("Content-Type", h.ContentType)
	}
	if h.ContentLength != "" {
		dst.Set("Content-Length", h.ContentLength)
	}
	if h.ContentEncoding != "" {
		dst.Set("Content-Encoding", h.ContentEncoding)
	}
	if h.ContentVersion != "" {
		dst.Set("Content-Version", h.ContentVersion)
	}
	if h.ContentLocation != "" {
		dst.Set("Content-Location", h.ContentLocation)
	}
	if h.ExtraFields != nil {
		for k, v := range h.ExtraFields {
			for _, v2 := range v {
				dst.Add(k, v2)
			}
		}
	}
}

func (H *THttp) completHeader(req *http.Request) {
	H.applyHeaders(req.Header)
}
func (H *THttp) completAutorization(req *http.Request) error {
	// Usa variáveis locais para não mutar H.AuthorizationType / H.Authorization,
	// senão a próxima chamada de Send() não renovaria o token Auth2.
	effectiveType := H.AuthorizationType
	bearerToken := H.Authorization

	if effectiveType == AT_AutoDetect {
		if bearerToken != "" {
			effectiveType = AT_Bearer
		} else if H.UserName != "" && H.Password != "" {
			effectiveType = AT_Basic
		}
	}
	if effectiveType == AT_Auth2 {
		token, err := H.Auth2.GetToken()
		if err != nil {
			return fmt.Errorf("erro ao obter o token: %v", err)
		}
		effectiveType = AT_Bearer
		bearerToken = token
	}
	if effectiveType == AT_Bearer {
		if strings.Contains(strings.ToLower(bearerToken), "bearer") {
			req.Header.Set("Authorization", bearerToken)
		} else {
			req.Header.Set("Authorization", "Bearer "+bearerToken)
		}
	}
	if effectiveType == AT_Basic {
		auth := H.UserName + ":" + H.Password
		basic := base64.StdEncoding.EncodeToString([]byte(auth))
		H.Request.Header.Authorization = "Basic " + basic
		req.SetBasicAuth(H.UserName, H.Password)
	}
	return nil
}
func (H *THttp) completAutorizationSocket(req http.Header) error {
	// Usa variáveis locais para não mutar H.AuthorizationType / H.Authorization
	// — assim Conectar() pode ser chamado várias vezes e o auto-detect / refresh
	// de token Auth2 continuam funcionando em chamadas subsequentes.
	effectiveType := H.AuthorizationType
	bearerToken := H.Authorization

	if effectiveType == AT_AutoDetect {
		if bearerToken != "" {
			effectiveType = AT_Bearer
		} else if H.UserName != "" && H.Password != "" {
			effectiveType = AT_Basic
		}
	}
	if effectiveType == AT_Auth2 {
		token, err := H.Auth2.GetToken()
		if err != nil {
			return fmt.Errorf("erro ao obter o token: %v", err)
		}
		effectiveType = AT_Bearer
		bearerToken = token
	}
	if effectiveType == AT_Bearer {
		if strings.Contains(strings.ToLower(bearerToken), "bearer") {
			req.Set("Authorization", bearerToken)
		} else {
			req.Set("Authorization", "Bearer "+bearerToken)
		}
	}
	if effectiveType == AT_Basic {
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

func (H *THttp) UseCert() bool {
	return (H.Certificate.PathCrt != "" && H.Certificate.PathPriv != "") ||
		(len(H.Certificate.CertPEMBlock) > 0 && len(H.Certificate.KeyPEMBlock) > 0) ||
		(len(H.Certificate.PfxBlock) > 0 && H.Certificate.PfxPass != "")

}

func (H *THttp) LoadPfx() (tls.Certificate, error) {
	privateKey, certificate, errPfx := pkcs12.Decode(H.Certificate.PfxBlock, H.Certificate.PfxPass)
	if errPfx != nil {
		return tls.Certificate{}, errPfx
	}
	cert := tls.Certificate{
		Certificate: [][]byte{certificate.Raw},
		PrivateKey:  privateKey,
		Leaf:        certificate,
	}
	return cert, nil
}

func (H *THttp) Send() (RES *Response, err error) {
	// Recover intencional: THttp é usado em muitos pontos da aplicação e
	// uma falha em uma chamada não pode derrubar as demais operações.
	// O panic é convertido em erro para o caller registrar em log e tratar
	// posteriormente — não remova sem alinhar antes.
	defer func() {
		if r := recover(); r != nil {
			RES = nil
			err = fmt.Errorf("recuperado de um panic no método Send: %v", r)
		}
	}()
	H.Response = NewResponse()

	var resp *http.Response
	var trans *http.Transport
	var Config *tls.Config
	var cert tls.Certificate

	// Obter transport (só será criado se necessário)
	trans = H.GetTransport()

	// Configurar certificados se necessário
	if H.UseCert() {
		if H.Certificate.PathCrt != "" && H.Certificate.PathPriv != "" {
			cert, err = tls.LoadX509KeyPair(H.Certificate.PathCrt, H.Certificate.PathPriv)
		} else if len(H.Certificate.CertPEMBlock) > 0 && len(H.Certificate.KeyPEMBlock) > 0 {
			cert, err = tls.X509KeyPair(H.Certificate.CertPEMBlock, H.Certificate.KeyPEMBlock)
		} else if len(H.Certificate.PfxBlock) > 0 && H.Certificate.PfxPass != "" {
			cert, err = H.LoadPfx()
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar certificado: %v", err)
		}
	}

	// Determinar se precisa de configuração TLS
	var needsTLS bool
	//var isHTTPS = strings.EqualFold(H.Protocolo, "HTTPS")

	// Precisa TLS se:
	needsTLS =
		H.TransportType == TTLS || // Forçado por TransportType
			H.TransportType == TSSLTLS || // Forçado por TransportType
			H.UseCert() ||
			H.InsecureSkipVerify // Tem certificados

	// Configurar TLS se necessário
	if needsTLS {
		Config = &tls.Config{
			InsecureSkipVerify: H.InsecureSkipVerify,
		}

		// Adicionar certificados se existirem
		if H.UseCert() {
			Config.Certificates = []tls.Certificate{cert}
			// Se for PFX, adicionar CA pool (opcional, só se quiser confiar na raiz do próprio certificado)
			if len(H.Certificate.PfxBlock) > 0 && H.Certificate.PfxPass != "" && cert.Leaf != nil {
				caPool := x509.NewCertPool()
				caPool.AddCert(cert.Leaf)
				Config.RootCAs = caPool
			}
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
	var req *http.Request
	switch H.EncType {
	case ET_NONE:
		req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, nil)

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
						return nil, fmt.Errorf("erro ao criar o arquivo %s: %v\n", v.FieldName, err)
					}
					_, err = fileWriter.Write([]byte(v.FieldValue))
					if err != nil {
						return nil, fmt.Errorf("erro ao escrever o arquivo %s: %v\n", v.FieldName, err)
					}
				} else {
					if err := multipartWriter.WriteField(v.FieldName, v.FieldValue); err != nil {
						return nil, fmt.Errorf("erro ao escrever campo %s: %v", v.FieldName, err)
					}
				}
			}
		}
		if H.Request.ItensContentText != nil {
			for _, v := range H.Request.ItensContentText {
				if err := multipartWriter.WriteField(v.Name, v.Value.Text()); err != nil {
					return nil, fmt.Errorf("erro ao escrever campo %s: %v", v.Name, err)
				}
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
		// Close() aqui não fecha o buffer — apenas escreve o boundary final do multipart.
		if err := multipartWriter.Close(); err != nil {
			return nil, fmt.Errorf("erro ao finalizar multipart: %v", err)
		}
		H.Request.Header.ContentType = multipartWriter.FormDataContentType()
		req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, &requestBody)

		// Defina o cabeçalho da requisição para indicar que está enviando dados com o formato multipart/form-data

	case ET_X_WWW_FORM_URLENCODED:
		//fmt.Println("CT_X_WWW_FORM_URLENCODED:")
		formData := url.Values{}
		if H.Request.ItensFormField != nil {
			for _, v := range H.Request.ItensFormField {
				formData.Add(v.FieldName, v.FieldValue)
			}
		}
		req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, strings.NewReader(formData.Encode()))
	case ET_RAW:
		req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, bytes.NewReader(H.Request.Body))
	case ET_BINARY:
		//fmt.Println("CT_BINARY:")
		fileBuffer := &bytes.Buffer{}
		fileBuffer.Reset()
		if H.Request.ItensContentBin != nil {
			for _, v := range H.Request.ItensContentBin {
				_, err := fileBuffer.Write(v.Value)
				if err != nil {
					return nil, fmt.Errorf("erro ao copiar os dados do arquivo %s para o buffer: %v", v.FileName, err)
				}
			}
		}
		req, err = http.NewRequest(GetMethodStr(H.Metodo), uri, fileBuffer)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição %s: %v", GetMethodStr(H.Metodo), err)
	}
	// Nota: o auto-detect (AT_AutoDetect → Bearer/Basic) acontece dentro de
	// completAutorization a partir de variáveis locais, sem mutar H.AuthorizationType.
	if H.AuthorizationType == AT_Auth2 && H.Auth2.AuthUrl != "" && (H.Auth2.AuthUrl == H.GetUrl() || H.GetUrl() == "") {
		RES, err = H.Auth2.Send()
		if err != nil {
			return nil, fmt.Errorf("erro ao fazer a requisição %s: %v", GetMethodStr(H.Metodo), err)
		}
		H.Response = RES
		return RES, nil

	} else {
		if err = H.completAutorization(req); err != nil {
			return nil, err
		}
		H.completHeader(req)
		resp, err = client.Do(req)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer a requisição %s: %v", GetMethodStr(H.Metodo), err)
	}
	// Guarda intencional: pelo contrato do net/http, client.Do não devolve
	// (nil, nil), mas em produção já foram observados casos (combinações de
	// transport/proxy/cert) em que ambos vêm nil. Sem este return o
	// `defer resp.Body.Close()` abaixo dispararia panic. Não remova.
	if resp == nil && err == nil {
		return nil, errors.New("retorno vazio, sem erro, possivelmente pode ser certificado invalido, configurei modo inseguro")
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
	headers := make(http.Header)
	if err := H.completAutorizationSocket(headers); err != nil {
		return err
	}
	H.applyHeaders(headers)

	dialer := websocket.DefaultDialer
	var (
		err  error
		err2 error
	)

	conn, _, err := dialer.Dial(H.GetUrl(), headers)
	if err != nil {
		if H.OnSend != nil {
			H.OnSend.Error("Erro na conexão: " + err.Error())
		} else {
			fmt.Printf("Erro na conexão: %v\n", err)
		}
		return fmt.Errorf("Erro na conexão: %v", err)
	}
	H.setConn(conn)
	H.setStatus(OPEN)
	if H.OnSend != nil {
		H.OnSend.Msg(MSG_CONECTADO)
	} else {
		fmt.Println(MSG_CONECTADO)
	}

	go func() {
		for {
			ws := H.getConn()

			autoReconnect := H.getAutoReconnect()
			if (ws == nil || err2 != nil) && autoReconnect {
				H.setStatus(CONNECTING)
				if H.getAttempts() >= H.getMaxAttempts() {
					break
				}
				if H.OnSend != nil {
					H.OnSend.Disconect(MSG_DISCONECT, false)
					H.OnSend.Msg(MSG_RECONECTANDO)
				} else {
					fmt.Printf(MSG_RECONECTANDO)
				}
				newConn, _, dialErr := dialer.Dial(H.GetUrl(), headers)
				if dialErr != nil {
					if H.OnSend != nil {
						H.OnSend.Error("Erro na conexão: " + dialErr.Error())
					} else {
						fmt.Printf("Erro na conexão: %v\n", dialErr)
					}
					time.Sleep(5 * time.Second)
					H.incAttempts()
					continue
				}
				H.setConn(newConn)
				H.setStatus(OPEN)
				H.resetAttempts()
				ws = newConn
				if H.OnSend != nil {
					H.OnSend.Msg(MSG_RECONECTADO)
				} else {
					fmt.Printf(MSG_RECONECTADO)
				}
			} else if (ws == nil || err2 != nil) && !autoReconnect {
				H.setStatus(CLOSED)
				break
			}
			err2 = nil

			for err2 == nil {
				ws = H.getConn()
				if ws == nil {
					err2 = errors.New("conexão WebSocket fechada")
					continue
				}
				msgtype, msg, readErr := ws.ReadMessage()
				if readErr != nil {
					if H.OnSend != nil {
						H.OnSend.Error("Erro na leitura da mensagem: " + readErr.Error())
					} else {
						fmt.Printf("Erro na leitura da mensagem: %v\n", readErr)
					}
					ws.Close()
					time.Sleep(5 * time.Second)
					err2 = readErr
					continue
				}
				if H.OnSend != nil {
					H.OnSend.Read(msgtype, msg, readErr)
				}
			}
		}
		if H.OnSend != nil {
			H.OnSend.Disconect(MSG_DISCONECT, true)
		} else {
			fmt.Printf("Erro na leitura da mensagem: %v\n", err)
		}
		H.setStatus(CLOSED)
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
			fmt.Printf("Erro na conexão: %v\n", err)
			fmt.Println("Tentando reconectar em 5 segundos...")
		}
		return err
		//time.Sleep(5 * time.Second)
	}
	//}
	return nil
}

// getConn devolve um snapshot do *websocket.Conn atual sob lock de leitura.
// Use o ponteiro retornado em vez de H.ws diretamente para evitar race
// com Free()/Desconectar()/reconexão na goroutine de leitura.
func (H *THttp) getConn() *websocket.Conn {
	H.wsMu.RLock()
	defer H.wsMu.RUnlock()
	return H.ws
}

func (H *THttp) setConn(c *websocket.Conn) {
	H.wsMu.Lock()
	H.ws = c
	H.wsMu.Unlock()
}

// setStatus / incAttempts / resetAttempts encapsulam mutações em
// H.WebSocket feitas pela goroutine de leitura.
func (H *THttp) setStatus(s Status) {
	H.wsMu.Lock()
	if H.WebSocket != nil {
		H.WebSocket.connect = s
	}
	H.wsMu.Unlock()
}

func (H *THttp) incAttempts() {
	H.wsMu.Lock()
	if H.WebSocket != nil {
		H.WebSocket.attempts++
	}
	H.wsMu.Unlock()
}

func (H *THttp) resetAttempts() {
	H.wsMu.Lock()
	if H.WebSocket != nil {
		H.WebSocket.attempts = 0
	}
	H.wsMu.Unlock()
}

func (H *THttp) getAttempts() int {
	H.wsMu.RLock()
	defer H.wsMu.RUnlock()
	if H.WebSocket == nil {
		return 0
	}
	return H.WebSocket.attempts
}

func (H *THttp) getAutoReconnect() bool {
	H.wsMu.RLock()
	defer H.wsMu.RUnlock()
	if H.WebSocket == nil {
		return false
	}
	return H.WebSocket.AutoReconnect
}

func (H *THttp) getMaxAttempts() int {
	H.wsMu.RLock()
	defer H.wsMu.RUnlock()
	if H.WebSocket == nil {
		return 0
	}
	return H.WebSocket.NumberOfAttempts
}

func (H *THttp) IsConect() bool {
	return H.getConn() != nil
}

func (H *THttp) Desconectar() error {
	H.wsMu.Lock()
	ws := H.ws
	H.ws = nil
	H.wsMu.Unlock()
	if ws == nil {
		return fmt.Errorf("conexão WebSocket não inicializada")
	}
	return ws.Close()
}

func (H *THttp) EnviarBinario(messageType int, data []byte) error {
	ws := H.getConn()
	if ws == nil {
		return fmt.Errorf("erro ao enviar mensagem, conexão não estabelecida")
	}
	return ws.WriteMessage(messageType, data)
}

func (H *THttp) EnviarTexto(messageType int, data string) error {
	ws := H.getConn()
	if ws == nil {
		return fmt.Errorf("erro ao enviar mensagem, conexão não estabelecida")
	}
	return ws.WriteMessage(messageType, []byte(data))
}

func (H *THttp) EnviarTextTypeTextMessage(data []byte) error {
	ws := H.getConn()
	if ws == nil {
		return fmt.Errorf("erro ao enviar mensagem, conexão não estabelecida")
	}
	return ws.WriteMessage(websocket.TextMessage, data)
}

func (H *THttp) EnviarBinarioTypeBinaryMessage(data []byte) error {
	ws := H.getConn()
	if ws == nil {
		return fmt.Errorf("erro ao enviar mensagem, conexão não estabelecida")
	}
	return ws.WriteMessage(websocket.BinaryMessage, data)
}
func (H *THttp) ConvertBodyInStruct(value any) error {
	err := json.Unmarshal(H.Response.Body, value)
	if err != nil {
		return err
	}
	return nil
}
