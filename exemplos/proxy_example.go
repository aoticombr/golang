package main

import (
	"fmt"

	http "github.com/aoticombr/golang/http"
)

func main() {
	// Exemplo 1: Requisição simples sem proxy
	fmt.Println("=== Requisição sem proxy ===")
	client1 := http.NewHttp()
	defer client1.Free()

	client1.SetUrl("https://httpbin.org/get")
	client1.SetMetodo(http.M_GET)

	fmt.Printf("Proxy ativo: %t\n", client1.GetProxyAtivo())

	// Exemplo 2: Configurar proxy
	fmt.Println("\n=== Configurando proxy ===")
	client2 := http.NewHttp()
	defer client2.Free()

	// Configurar proxy básico
	client2.SetProxyConfig("proxy.empresa.com", 8080, "usuario", "senha")

	fmt.Printf("Proxy ativo: %t\n", client2.GetProxyAtivo())
	fmt.Printf("Proxy host: %s\n", client2.GetProxyHost())
	fmt.Printf("Proxy port: %d\n", client2.GetProxyPort())
	fmt.Printf("Proxy user: %s\n", client2.GetProxyUserName())

	// Exemplo 3: Configurar propriedades individuais
	fmt.Println("\n=== Configuração individual ===")
	client3 := http.NewHttp()
	defer client3.Free()

	client3.SetProxyHost("meu-proxy.com")
	client3.SetProxyPort(3128)
	client3.SetProxyUserName("meuusuario")
	client3.SetProxyPassword("minhasenha")
	client3.SetProxyAtivo(true)

	fmt.Printf("Proxy configurado: %s:%d\n", client3.GetProxyHost(), client3.GetProxyPort())
	fmt.Printf("Usuário: %s\n", client3.GetProxyUserName())
	fmt.Printf("Proxy ativo: %t\n", client3.GetProxyAtivo())

	// Desabilitar proxy
	client3.SetProxyAtivo(false)
	fmt.Printf("Proxy após desabilitar: %t\n", client3.GetProxyAtivo())

	// Exemplo 4: Mostrando que transport só é criado quando necessário
	fmt.Println("\n=== Testando criação de transport ===")

	// Cliente sem nada especial - não deve criar transport
	clientSimples := http.NewHttp()
	defer clientSimples.Free()
	transport1 := clientSimples.GetTransport()
	fmt.Printf("Transport sem configurações especiais: %v\n", transport1)

	// Cliente com proxy - deve criar transport
	clientComProxy := http.NewHttp()
	defer clientComProxy.Free()
	clientComProxy.SetProxyConfig("proxy.test.com", 8080, "", "")
	transport2 := clientComProxy.GetTransport()
	fmt.Printf("Transport com proxy: %v\n", transport2 != nil)
}
