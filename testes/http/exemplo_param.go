package main

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://127.0.0.1:3003/signin?eee=1111&aaaa=222222&bbbbbbbbb=3333333")
	fmt.Println("URL:", cp.GetUrl())
	cp.Params.Add("teste", "teste")
	cp.Params.Set("aaaa", "999999")
	for k, v := range cp.Params {
		fmt.Println("Params:", k, v)
	}
	fmt.Println("URL:", cp.GetUrl())

}
