package main

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	//	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://127.0.0.1:3003/signin/{{id}}")
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
