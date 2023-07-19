package main

import (
	"fmt"
	"io/ioutil"
	"os"

	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.Url = "http://127.0.0.1:3003/signin"

	cp.Metodo = http.M_POST
	//cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.ContentType = "application/octet-stream"
	cp.Request.Header.Accept = "*/*"
	cp.Request.AddContentBin("file1", "file1.txt", []byte("teste"))
	file, err := os.Open("Mickey_Mouse.png") // Substitua pelo caminho real do arquivo que deseja enviar
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()
	// Ler o conteúdo do arquivo como um slice de bytes
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Erro ao ler o conteúdo do arquivo:", err)
		return
	}
	cp.Request.AddContentBin("file2", "Mickey_Mouse.png", fileContent)

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
