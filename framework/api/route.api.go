package api

import (
	"fmt"
	"net/http"
)

type TipoMetodoEnum int

const (
	GET TipoMetodoEnum = iota
	POST
	PUT
	DELETE
	PATCH
	WEBSOCKET
)

type Metodo struct {
	TipoMetodo TipoMetodoEnum
	Funcao     http.HandlerFunc
}

type Rotas map[string][]Metodo

func (r Rotas) RegisterMetodo(rota string, tipoMetodo TipoMetodoEnum, funcao http.HandlerFunc) {
	fmt.Println("Registrando rota: ", rota, tipoMetodo)

	r[rota] = append(r[rota], Metodo{tipoMetodo, funcao})
}

type RegistraRotasController struct {
	RegisteredRoutes Rotas
}

var registraRotasGlobal *RegistraRotasController

func GetRegistraRotasInstance() *RegistraRotasController {
	if registraRotasGlobal == nil {
		registraRotasGlobal = &RegistraRotasController{
			RegisteredRoutes: make(Rotas),
		}
	}
	return registraRotasGlobal
}

func RegistrarRota(rota string, tipoMetodo TipoMetodoEnum, funcao http.HandlerFunc) {
	GetRegistraRotasInstance().RegisteredRoutes.RegisterMetodo(rota, tipoMetodo, funcao)
}
