package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	http2 "github.com/aoticombr/golang/http"
)

func Teste(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Teste")
	cp := http2.NewHttp()
	cp.SetUrl("http://ws3007.aoti.com.br/DB1/v1/ping")
	///	cp.UserName = "thiago.silva@nbsi.com.br"
	///	cp.Password = "Paymail01@"
	cp.Metodo = http2.M_GET
	cp.EncType = http2.ET_RAW
	cp.Request.Header.ContentType = "application/json"
	resp, err := cp.Send()
	if err != nil {
		fmt.Println(err)
	}
	w.Write(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func MemoryLeak() {
	go func() {
		r := http.NewServeMux()
		// Define uma rota usando o Gorilla Mux
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Escreve uma resposta no corpo da resposta HTTP
			fmt.Fprintf(w, "Olá! Você acessou: %s\n", r.URL.Path)
		})
		r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		r.Handle("/debug/pprof/{cmd}", http.HandlerFunc(pprof.Index)) // special handling for Gorilla mux
		r.Handle("/debug/pprof/heap", http.HandlerFunc(pprof.Handler("heap").ServeHTTP))
		r.Handle("/debug/pprof/goroutine", http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP))
		r.Handle("/debug/pprof/block", http.HandlerFunc(pprof.Handler("block").ServeHTTP))
		r.Handle("/debug/pprof/allocs", http.HandlerFunc(pprof.Handler("allocs").ServeHTTP))
		r.Handle("/debug/pprof/threadcreate", http.HandlerFunc(pprof.Handler("threadcreate").ServeHTTP))
		r.Handle("/teste", http.HandlerFunc(Teste))

		// Inicia o servidor HTTP usando o roteador criado pelo Gorilla Mux
		fmt.Println("Servidor iniciado em: http://localhost:6060")
		if err := http.ListenAndServe(":6060", r); err != nil {
			fmt.Printf("Erro ao iniciar o servidor: %s\n", err)
		}
	}()
	//go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap
}

func main() {
	MemoryLeak()
	select {}
}
