package main

import (
	"time"
)

type ContasPagar struct {
	Cdcontaspagar int64
	Historico     string
	Dtaconta      time.Time
	Valor         float32
}

type any struct {
	// Defina os campos necessários para o tipo `any`
}

type value interface {
	GetData() interface{} // Exemplo de um método na interface `value`
}

type component struct {
	Value value
	// Defina outros campos necessários para o componente
}

func main() {

}
