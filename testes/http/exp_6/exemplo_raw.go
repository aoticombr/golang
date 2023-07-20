package main

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("https://api.xolvis.com/pm/oauth/token")
	cp.UserName = "thiago.silva@nbsi.com.br"
	cp.Password = "Paymail01@"
	cp.Metodo = http.M_POST
	cp.Request.Header.ContentType = "application/x-www-form-urlencoded"
	cp.Request.AddFormField("grant_type", "client_credentials")
	// cp.Request.Body = []byte(`{
	// 	"user":"admin@aoti.com.br",
	// 	"pass":"master"
	// }	`)
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
