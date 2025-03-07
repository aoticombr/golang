package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/aoticombr/golang/component"
	"github.com/aoticombr/golang/http"
)

func TestAuth2_tipo1(t *testing.T) {
	fmt.Println("Teste")
	cp1 := http.NewHttp()
	cp1.SetUrl("http://100.0.66.81:3003/token3")
	cp1.AuthorizationType = http.AT_Auth2
	cp1.Auth2.ClientId = "ddddddddd"
	cp1.Auth2.ClientSecret = "fffffff"
	cp1.Auth2.AuthUrl = "http://100.0.66.81:3003/token"
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
	cp2.Auth2.AuthUrl = "http://localhost:3003/token"
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
	a := component.NewStrings()
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
