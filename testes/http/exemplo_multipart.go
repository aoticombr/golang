package main

import (
	"fmt"
	"io/ioutil"
	"os"

	comp "github.com/aoticombr/golang/component"
	http "github.com/aoticombr/golang/http"
)

func main() {
	fmt.Println("Teste")
	cp := http.NewHttp()
	cp.SetUrl("http://127.0.0.1:3003/signin?eee=1111&aaaa=222222&bbbbbbbbb=3333333")

	cp.Metodo = http.M_POST
	//cp.Request.Header.ContentType = "application/json"
	cp.Request.Header.ContentType = "multipart/form-data"
	cp.Request.Header.ContentType = "multipart/form-data"
	cp.Request.Header.Accept = "*/*"
	t := comp.NewStrings().Add("xxxxxyyyy").Add("eeeeee")
	cp.Request.AddContentText("txt1", t)
	cp.Request.AddContentBin("file1", "file1.txt", []byte("teste"))
	file, err := os.Open("image.png") // Substitua pelo caminho real do arquivo que deseja enviar
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
	cp.Request.AddContentBin("file2", "image.png", fileContent)

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
