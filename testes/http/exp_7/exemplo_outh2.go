package main

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("https://...")
	cp.AuthorizationType = http.AT_Auth2
	cp.Auth2.ClientId = ".."
	cp.Auth2.ClientSecret = "....."
	cp.Auth2.AuthUrl = "...."
	cp.Metodo = http.M_GET
	cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.AddField("X-Personal-ID", "...")
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
