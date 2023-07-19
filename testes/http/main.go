package main

import (
	"fmt"
	"io/ioutil"
	"os"

	comp "github.com/aoticombr/golang/component"
	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.Url = "http://127.0.0.1:3003/signin"

	cp.Metodo = http.M_POST
	//cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.ContentType = "multipart/form-data"
	cp.Request.Header.Accept = "*/*"
	cp.Request.Header.AcceptCharset = "utf-8"
	cp.Request.Header.AcceptEncoding = "gzip, deflate, br"
	cp.Request.Header.AcceptLanguage = "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7"
	cp.Request.Header.Authorization = "Bearer teste"
	cp.Request.Header.Charset = "utf-8"
	cp.Request.Header.ContentLocation = "http://"
	//cp.Request.Header.ContentLength = "0"
	//cp.Request.Header.ContentEncoding = "gzip"
	//cp.Request.Header.ContentVersion = "1.0"

	cp.Request.AddFormField("teste", "teste")
	cp.Request.AddFormField("teste2", "teste2")
	t := comp.NewStrings().Add("xxxxxyyyy").Add("eeeeee")
	cp.Request.AddContentText("txt1", t)
	cp.Request.AddContentBin("file1", "file1.txt", []byte("teste"))
	file, err := os.Open("12087033_898785803548607_2614616143038690718_o.jpg") // Substitua pelo caminho real do arquivo que deseja enviar
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
	cp.Request.AddContentBin("file2", "12087033_898785803548607_2614616143038690718_o.jpg", fileContent)

	cp.Request.Header.AddExtraField("testexx", "testexx")
	cp.Request.Header.AddExtraField("testexx1", "testexx1")
	cp.Request.AddSubmitFile("teste", "application/json", []byte("teste"))

	cp.Request.Body = []byte(`{
		"user":"admin@aoti.com.br",
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
