package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aoticombr/golang/http"
	"github.com/aoticombr/golang/stringlist"
)

func TestAuth2_tipo1(t *testing.T) {
	fmt.Println("Teste")
	cp1 := http.NewHttp()
	cp1.SetUrl("http://100.0.66.81:3003/token3")
	cp1.AuthorizationType = http.AT_Auth2
	cp1.Auth2.ClientId = "ddddddddd"
	cp1.Auth2.ClientSecret = "fffffff"
	cp1.Auth2.AccessTokenUrl = "http://100.0.66.81:3003/token"
	cp1.Auth2.Scope = "downloaded"
	cp1.Auth2.ClientAuth = http.CA_SendBasicAuthHeader
	cp1.Metodo = http.M_GET
	cp1.EncType = http.ET_RAW
	cp1.Request.Header.ContentType = "application/json"
	///cp1.Request.Header.AddField("X-Personal-ID", "...")
	cp1.Request.Body = []byte(``)
	resp, err := cp1.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	for k, v := range resp.Header {
		fmt.Println("Header:", k, v)
	}
	fmt.Println("Body:", resp.Body)
	fmt.Println("Body string:", string(resp.Body))

}
func TestAuth2_tipo2(t *testing.T) {
	fmt.Println("Teste")
	cp2 := http.NewHttp()
	cp2.SetUrl("http://localhost:3003/token3")
	cp2.AuthorizationType = http.AT_Auth2
	cp2.Auth2.ClientId = "ddddddddd"
	cp2.Auth2.ClientSecret = "fffffff"
	cp2.Auth2.AccessTokenUrl = "http://localhost:3003/token"
	cp2.Auth2.Scope = "downloaded"
	cp2.Auth2.ClientAuth = http.CA_SendClientCredentialsInBody
	cp2.Metodo = http.M_GET
	cp2.EncType = http.ET_RAW
	cp2.Request.Header.ContentType = "application/json"
	///cp2.Request.Header.AddField("X-Personal-ID", "...")
	cp2.Request.Body = []byte(``)
	resp, err := cp2.Send()
	if err != nil {
		fmt.Println("Erro:", err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	for k, v := range resp.Header {
		fmt.Println("Header:", k, v)
	}
	fmt.Println("Body:", resp.Body)
	fmt.Println("Body string:", string(resp.Body))

}
func TestGetRaw(t *testing.T) {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("https://aoti123.free.beeceptor.com")

	cp.Metodo = http.M_GET
	cp.EncType = http.ET_RAW
	cp.Request.Header.ContentType = "application/json"
	resp, err := cp.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	for k, v := range resp.Header {
		fmt.Println("Header:", k, v)
	}
	fmt.Println("Body:", resp.Body)
	fmt.Println("Body string:", string(resp.Body))
}
func TestSendRaw(t *testing.T) {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://127.0.0.1:3003")

	cp.Metodo = http.M_POST
	cp.EncType = http.ET_RAW
	cp.Request.Header.ContentType = "application/json"
	cp.Request.AddFormField("grant_type", "client_credentials")
	cp.Request.Body = []byte(`{
	 	"user":"admin@teste.com.br",
	 	"pass":"master"
	 }	`)
	resp, err := cp.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	for k, v := range resp.Header {
		fmt.Println("Header:", k, v)
	}
	fmt.Println("Body:", resp.Body)
	fmt.Println("Body string:", string(resp.Body))
}

func TestSendParam(t *testing.T) {
	//	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.EncType = http.ET_RAW
	cp.SetUrl("http://127.0.0.1:3003/{{id}}")
	//	fmt.Println("Path:", cp.Path)
	//	fmt.Println("URL:", cp.GetUrl())
	cp.Params.Add("teste", "teste")
	cp.Params.Set("aaaa", "999999")
	//cp.Varibles.Add("id", "123456789")
	// for k, v := range cp.Params {
	// 	fmt.Println("Params:", k, v)
	// }
	// for k, v := range cp.Varibles {
	// 	fmt.Println("Varibles:", k, v)
	// }
	//	fmt.Println("URL:", cp.GetUrl())
	resp, err := cp.Send()
	if err != nil {
		fmt.Println("Erro:", err)
	}
	fmt.Println("Status:", resp)
}

func TestMultPart(t *testing.T) {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://localhost:3003/?eee=1111&aaaa=222222&bbbbbbbbb=3333333")

	cp.Metodo = http.M_POST
	cp.EncType = http.ET_FORM_DATA
	cp.Request.Header.ContentType = "multipart/form-data"
	a := stringlist.NewStrings()
	a.Add("xxxxxyyyy")
	a.Add("eeeeee")
	cp.Request.AddContentText("txt1", a)
	file, err := os.Open("image.png") // Substitua pelo caminho real do arquivo que deseja enviar
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()
	// Ler o conteúdo do arquivo como um slice de bytes
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Erro ao ler o conteúdo do arquivo:", err)
		return
	}
	cp.Request.AddContentBin("file2", "image.png", fileContent)

	resp, err := cp.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	for k, v := range resp.Header {
		fmt.Println("Header:", k, v)
	}
	fmt.Println("Body:", resp.Body)
	fmt.Println("Body string:", string(resp.Body))
}

func TestFormData(t *testing.T) {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://localhost:3003")

	cp.Metodo = http.M_POST
	cp.EncType = http.ET_X_WWW_FORM_URLENCODED
	//cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.ContentType = "application/x-www-form-urlencoded"

	cp.Request.AddFormField("teste", "teste")
	cp.Request.AddFormField("teste2", "teste2")

	cp.Request.Header.AddField("testexx", "testexx")
	cp.Request.Header.AddField("testexx1", "testexx1")

	resp, err := cp.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	for k, v := range resp.Header {
		fmt.Println("Header:", k, v)
	}
	fmt.Println("Body:", resp.Body)
	fmt.Println("Body string:", string(resp.Body))
}

func TestBinary(t *testing.T) {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://localhost:3003")

	cp.Metodo = http.M_POST
	cp.EncType = http.ET_BINARY
	//cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.ContentType = "application/octet-stream"
	file, err := os.Open("image.png") // Substitua pelo caminho real do arquivo que deseja enviar
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()
	// Ler o conteúdo do arquivo como um slice de bytes
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Erro ao ler o conteúdo do arquivo:", err)
		return
	}
	cp.Request.Body = []byte(`{}`)
	cp.Request.AddContentBin("file2", "image.png", fileContent)

	resp, err := cp.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	for k, v := range resp.Header {
		fmt.Println("Header:", k, v)
	}
	//fmt.Println("Body:", resp.Body)
	//fmt.Println("Body string:", string(resp.Body))
}

func TestBinaryType(t *testing.T) {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://100.0.66.81:3003")

	cp.Metodo = http.M_POST
	cp.EncType = http.ET_FORM_DATA
	//cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.ContentType = "multipar/form-data"
	file, err := os.Open("H:\\golang\\testes\\http\\testepaulo.pdf") // Substitua pelo caminho real do arquivo que deseja enviar
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()
	// Ler o conteúdo do arquivo como um slice de bytes
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Erro ao ler o conteúdo do arquivo:", err)
		return
	}
	//cp.Request.Body = []byte(`{}`)
	cp.Request.AddSubmitFile("", "testepaulo.pdf", "application/pdf", fileContent)

	resp, err := cp.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	for k, v := range resp.Header {
		fmt.Println("Header:", k, v)
	}
	//fmt.Println("Body:", resp.Body)
	//fmt.Println("Body string:", string(resp.Body))
}

type ReadSocket struct {
}

func (rs *ReadSocket) Read(messageType int, body []byte, err error) {
	fmt.Println("--------------")
	fmt.Println("ReadSocket.read")
	fmt.Println("messageType:", messageType)
	fmt.Println("body:", string(body))
	fmt.Println("err:", err)
}
func (rs *ReadSocket) Error(msg string) {
	fmt.Println("-------Error-------")
	fmt.Println(time.Now())
	fmt.Println("msg:", msg)
}
func (rs *ReadSocket) Msg(msg string) {
	fmt.Println("-------Msg-------")
	fmt.Println(time.Now())
	fmt.Println("msg:", msg)
}
func (rs *ReadSocket) Disconect(msg string, limit bool) {
	fmt.Println("-------Disconect-------")
	fmt.Println(time.Now())
	fmt.Println("msg:", msg)
	fmt.Println("limit:", limit)
}

func TestWebSocket(t *testing.T) {
	var rs *ReadSocket
	rs = &ReadSocket{}
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.Request.Header.AddField("x-xxx-dealer", "07600973")
	cp.Authorization = `eyJg`
	cp.AuthorizationType = http.AT_Bearer
	cp.SetUrl("ws://localhost:3030/route1")

	cp.Metodo = http.M_POST
	cp.EncType = http.ET_WEB_SERVICE

	cp.OnSend = rs
	fmt.Println("111111")
	err := cp.Conectar()
	if err != nil {
		panic(err)
	}
	fmt.Println("222222")
	go func() {
		for {
			fmt.Println("44444")
			err = cp.EnviarTextTypeTextMessage([]byte("Teste"))
			if err != nil {
				fmt.Println("Erro ao enviar:", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	fmt.Println("33333333")
	select {}
}

// TestAuth2_AuthorizationCode_Url monta a URL de autorização que deve ser aberta no
// navegador. Abra a URL impressa, autentique e o IdP redireciona para o CallBackUrl
// com ?code=... — copie esse code para usar em TestAuth2_AuthorizationCode_Token.
func TestAuth2_AuthorizationCode_Url(t *testing.T) {
	cp := http.NewHttp()
	cp.AuthorizationType = http.AT_Auth2
	cp.Auth2.GrantType = http.GT_AuthorizationCode
	cp.Auth2.AuthUrl = "https://login.microsoftonline.com/{tenantId}/oauth2/v2.0/authorize"
	cp.Auth2.AccessTokenUrl = "https://login.microsoftonline.com/{tenantId}/oauth2/v2.0/token"
	cp.Auth2.CallBackUrl = ""
	cp.Auth2.ClientId = ""
	cp.Auth2.ClientSecret = ""
	cp.Auth2.Scope = ""

	loginURL, err := cp.Auth2.GetAuthorizationUrl("teste-state-123")
	if err != nil {
		t.Fatalf("erro ao montar a URL de autorização: %v", err)
	}

	fmt.Println("Abra esta URL no navegador para obter o code:")
	fmt.Println(loginURL)
}

// TestAuth2_AuthorizationCode_Token troca o authorization code (obtido no redirect do
// navegador via TestAuth2_AuthorizationCode_Url) por um access token.
//
// O code é de uso único: preencha `code` com o valor capturado no CallBackUrl e rode o
// teste. Sem o code o teste é pulado (não há como obtê-lo de forma automatizada).
func TestAuth2_AuthorizationCode_Token(t *testing.T) {
	const code = ""

	if code == "" {
		t.Skip("preencha a const `code` com o valor retornado no CallBackUrl para rodar este teste")
	}

	cp := http.NewHttp()
	cp.AuthorizationType = http.AT_Auth2
	cp.Auth2.GrantType = http.GT_AuthorizationCode
	cp.Auth2.AuthUrl = "https://login.microsoftonline.com/{tenantId}/oauth2/v2.0/authorize"
	cp.Auth2.AccessTokenUrl = "https://login.microsoftonline.com/{tenantId}/oauth2/v2.0/token"
	cp.Auth2.CallBackUrl = ""
	cp.Auth2.ClientId = ""
	cp.Auth2.ClientSecret = ""
	cp.Auth2.Scope = ""
	// Azure AD envia client_id/secret no corpo no fluxo authorization_code.
	cp.Auth2.ClientAuth = http.CA_SendClientCredentialsInBody
	cp.Auth2.Code = code

	const tenantID = ""

	token, err := cp.Auth2.GetToken()
	if err != nil {
		t.Fatalf("erro ao trocar o code por token: %v", err)
	}

	// 1) HTTP do endpoint de token deve ter sido 2xx
	if cp.Auth2.Resp == nil {
		t.Fatal("resposta do endpoint de token é nil")
	}
	if cp.Auth2.Resp.StatusCode < 200 || cp.Auth2.Resp.StatusCode >= 300 {
		t.Fatalf("status inesperado do endpoint de token: %d %s\nBody: %s",
			cp.Auth2.Resp.StatusCode, cp.Auth2.Resp.StatusMessage, string(cp.Auth2.Resp.Body))
	}

	// 2) token não pode ser vazio
	if token == "" {
		t.Fatal("token vazio retornado")
	}

	// 3) valida que é um JWT do tenant esperado, não expirado e com a audience certa
	claims, err := decodeJWTClaims(token)
	if err != nil {
		t.Fatalf("access token não é um JWT válido: %v\nToken: %s", err, token)
	}

	if iss, _ := claims["iss"].(string); !strings.Contains(iss, tenantID) {
		t.Errorf("iss não pertence ao tenant esperado: iss=%q (esperava conter %q)", iss, tenantID)
	}
	if exp, ok := claims["exp"].(float64); !ok {
		t.Error("claim exp ausente no token")
	} else if int64(exp) <= time.Now().Unix() {
		t.Errorf("token já expirado: exp=%d, agora=%d", int64(exp), time.Now().Unix())
	}
	if appid, _ := claims["appid"].(string); appid != "" && appid != cp.Auth2.ClientId {
		t.Errorf("appid do token (%q) difere do ClientId usado (%q)", appid, cp.Auth2.ClientId)
	}

	fmt.Println("Access token OK")
	fmt.Println("  iss:", claims["iss"])
	fmt.Println("  aud:", claims["aud"])
	fmt.Println("  appid:", claims["appid"])
	fmt.Println("  scp:", claims["scp"])
	fmt.Println("  exp:", claims["exp"])
	fmt.Println("Token:", token)
}

func TestAuth2_AuthorizationWithPCKCECode_Token(t *testing.T) {
	const code = "" // <-- cole aqui o code retornado em ?code=... no CallBackUrl

	if code == "" {
		t.Skip("preencha a const `code` com o valor retornado no CallBackUrl para rodar este teste")
	}

	cp := http.NewHttp()
	cp.AuthorizationType = http.AT_Auth2
	cp.Auth2.GrantType = http.GT_AuthorizationCodeWithPKCE
	cp.Auth2.AuthUrl = "https://login.microsoftonline.com/{tenantID}/oauth2/v2.0/authorize"
	cp.Auth2.AccessTokenUrl = "https://login.microsoftonline.com/{tenantID}/oauth2/v2.0/token"
	//cp.Auth2.CallBackUrl = "http://localhost:8080/callback"
	cp.Auth2.ClientId = ""     //85db....
	cp.Auth2.ClientSecret = "" //Cds.....
	cp.Auth2.Scope = ""        //mail.read
	// Azure AD envia client_id/secret no corpo no fluxo authorization_code.
	cp.Auth2.ClientAuth = http.CA_SendClientCredentialsInBody
	cp.Auth2.Code = code

	//	PKCE (RFC 7636) - usado quando GrantType = GT_AuthorizationCodeWithPKCE.
	cp.Auth2.CodeVerifier = ""
	//	CodeChallengeMethod é “S256” (padrão) ou “plain”
	cp.Auth2.CodeChallengeMethod = "S256"

	const tenantID = ""

	token, err := cp.Auth2.GetToken()
	if err != nil {
		t.Fatalf("erro ao trocar o code por token: %v", err)
	}

	// 1) HTTP do endpoint de token deve ter sido 2xx
	if cp.Auth2.Resp == nil {
		t.Fatal("resposta do endpoint de token é nil")
	}
	if cp.Auth2.Resp.StatusCode < 200 || cp.Auth2.Resp.StatusCode >= 300 {
		t.Fatalf("status inesperado do endpoint de token: %d %s\nBody: %s",
			cp.Auth2.Resp.StatusCode, cp.Auth2.Resp.StatusMessage, string(cp.Auth2.Resp.Body))
	}

	// 2) token não pode ser vazio
	if token == "" {
		t.Fatal("token vazio retornado")
	}

	// 3) valida que é um JWT do tenant esperado, não expirado e com a audience certa
	claims, err := decodeJWTClaims(token)
	if err != nil {
		t.Fatalf("access token não é um JWT válido: %v\nToken: %s", err, token)
	}

	if iss, _ := claims["iss"].(string); !strings.Contains(iss, tenantID) {
		t.Errorf("iss não pertence ao tenant esperado: iss=%q (esperava conter %q)", iss, tenantID)
	}
	if exp, ok := claims["exp"].(float64); !ok {
		t.Error("claim exp ausente no token")
	} else if int64(exp) <= time.Now().Unix() {
		t.Errorf("token já expirado: exp=%d, agora=%d", int64(exp), time.Now().Unix())
	}
	if appid, _ := claims["appid"].(string); appid != "" && appid != cp.Auth2.ClientId {
		t.Errorf("appid do token (%q) difere do ClientId usado (%q)", appid, cp.Auth2.ClientId)
	}

	fmt.Println("Access token OK")
	fmt.Println("  iss:", claims["iss"])
	fmt.Println("  aud:", claims["aud"])
	fmt.Println("  appid:", claims["appid"])
	fmt.Println("  scp:", claims["scp"])
	fmt.Println("  exp:", claims["exp"])
	fmt.Println("Token:", token)
}

// decodeJWTClaims decodifica o payload (2ª parte) de um JWT sem validar a assinatura,
// retornando as claims como mapa. Serve para inspeção em teste, não para validação de segurança.
func decodeJWTClaims(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("formato JWT inválido: esperava 3 segmentos, obteve %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar payload base64url: %w", err)
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do payload: %w", err)
	}
	return claims, nil
}

// TestAws_SQS_ListQueues faz uma chamada real ao SQS (Query API) assinada com
// AWS Signature V4. Preencha AccessKey/SecretKey. Esperado: HTTP 200 com um XML
// <ListQueuesResponse>. Se as credenciais estiverem erradas, vem 403 com
// <Code>InvalidClientTokenId</Code> ou <Code>SignatureDoesNotMatch</Code>.
func TestAws_SQS_ListQueues(t *testing.T) {
	const (
		accessKey = "" // <-- preencha
		secretKey = "" // <-- preencha
	)
	if accessKey == "" || secretKey == "" {
		t.Skip("preencha accessKey/secretKey para rodar o teste real de AWS SigV4")
	}

	cp := http.NewHttp()
	cp.SetUrl("https://sqs.us-east-1.amazonaws.com/{id}/dev-crmi-lead-onecrm-nbs-queue.fifo?Action=ReceiveMessage&MaxNumberOfMessages=10&WaitTimeSeconds=5&AttributeName.1=All&MessageAttributeName.1=All")
	cp.SetMetodo(http.M_GET)
	cp.AuthorizationType = http.AT_AwsSignature
	cp.Aws.AccessKey = accessKey
	cp.Aws.SecretKey = secretKey
	cp.Aws.Region = "us-east-1"
	cp.Aws.Service = "sqs"
	// cp.Aws.SessionToken = "..." // se usar credenciais temporarias (STS)

	resp, err := cp.Send()
	if err != nil {
		t.Fatalf("erro ao enviar: %v", err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Body:", string(resp.Body))

	if resp.StatusCode != 200 {
		t.Errorf("esperava HTTP 200, veio %d: %s", resp.StatusCode, string(resp.Body))
	}
}
