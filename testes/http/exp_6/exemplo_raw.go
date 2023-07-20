package exp_6

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://127.0.0.1:3003/signin?eee=1111&aaaa=222222&bbbbbbbbb=3333333")

	cp.Metodo = http.M_POST
	cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.Accept = "*/*"
	cp.Request.Header.AcceptCharset = "utf-8"
	cp.Request.Header.AcceptEncoding = "gzip, deflate, br"
	cp.Request.Header.AcceptLanguage = "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7"
	cp.Request.Header.Authorization = "Bearer teste"
	cp.Request.Header.Charset = "utf-8"
	cp.Request.Header.ContentLocation = "http://"

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
