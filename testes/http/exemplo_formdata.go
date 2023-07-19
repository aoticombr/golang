package main

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.Url = "http://127.0.0.1:3003/signin"

	cp.Metodo = http.M_POST
	//cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.ContentType = "application/x-www-form-urlencoded"
	cp.Request.Header.Accept = "*/*"
	cp.Request.Header.AcceptCharset = "utf-8"
	cp.Request.Header.AcceptEncoding = "gzip, deflate, br"
	cp.Request.Header.AcceptLanguage = "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7"
	cp.Request.Header.Authorization = "Bearer teste"
	cp.Request.Header.Charset = "utf-8"
	cp.Request.Header.ContentLocation = "http://"

	cp.Request.AddFormField("teste", "teste")
	cp.Request.AddFormField("teste2", "teste2")

	cp.Request.Header.AddExtraField("testexx", "testexx")
	cp.Request.Header.AddExtraField("testexx1", "testexx1")
	cp.Request.AddSubmitFile("teste", "application/json", []byte("teste"))

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
