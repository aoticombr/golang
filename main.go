package main

import (
	"fmt"

	"github.com/aoticombr/golang/http"
)

func main() {
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
