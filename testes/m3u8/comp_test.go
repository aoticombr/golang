package main

import (
	"fmt"
	"testing"

	httpaoti "github.com/aoticombr/golang/http"
)

func TestUp_Down(t *testing.T) {
	fmt.Printf("teste")
	link2 := httpaoti.NewHttp()
	link2.SetUrl("https://url....")
	link2.SetMetodo(httpaoti.M_GET)
	link2.Request.Header.Authorization = "Bearer "
	link2.Request.Header.AddField("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36")
	link2.Request.Header.AddField("origin", "https://app..com.br")
	link2.Request.Header.AddField("referer", "https://app..com.br/")
	resp2, err := link2.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp2.Body))

}
