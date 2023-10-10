package main

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://localhost:3003")
	cp.AuthorizationType = http.AT_Auth2
	cp.Auth2.ClientId = "...."
	cp.Auth2.ClientSecret = "...."
	cp.Auth2.AuthUrl = "https://...."
	cp.Metodo = http.M_POST
	cp.Request.Header.ContentType = "application/json"
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
