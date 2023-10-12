package main

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://localhost:3003/token3")
	cp.AuthorizationType = http.AT_Auth2
	cp.Auth2.ClientId = "ddddddddd"
	cp.Auth2.ClientSecret = "fffffff"
	cp.Auth2.AuthUrl = "http://localhost:3003/token"
	cp.Auth2.Scope = "downloaded"
	cp.Auth2.ClientAuth = http.CA_SendClientCredentialsInBody
	cp.Metodo = http.M_GET
	cp.Request.Header.ContentType = "application/json"
	///cp.Request.Header.AddField("X-Personal-ID", "...")
	cp.Request.Body = []byte(``)
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
