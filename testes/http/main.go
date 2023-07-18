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
	resp, err := cp.Send()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("Msg:", resp.StatusMessage)
	fmt.Println("Header:", resp.Header)
	fmt.Println("Body:", resp.Body)
	fmt.Println("Body string:", string(resp.Body))
}
