package exp_3

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://127.0.0.1:3003/signin")

	cp.Metodo = http.M_POST
	//cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.ContentType = "application/x-www-form-urlencoded"

	cp.Request.AddFormField("teste", "teste")
	cp.Request.AddFormField("teste2", "teste2")

	cp.Request.Header.AddField("testexx", "testexx")
	cp.Request.Header.AddField("testexx1", "testexx1")

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