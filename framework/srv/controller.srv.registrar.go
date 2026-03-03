package srv

import (
	"fmt"
	"reflect"

	"github.com/aoticombr/golang/lib"
)

type RegisterClasses map[string]reflect.Type

type RegistraSrvController struct {
	RegisteredClasses RegisterClasses
}

var registraSrvGlobal *RegistraSrvController

func CriarStruct(nome string, mapaDeTipos map[string]reflect.Type) (interface{}, error) {
	// Verifica se o nome está presente no mapa
	tipo, ok := mapaDeTipos[nome]
	if !ok {
		return nil, fmt.Errorf("Struct com o nome %s não encontrada", nome)
	}

	// Cria uma instância da struct como um ponteiro
	instancia := reflect.New(tipo).Interface()

	return instancia, nil
}

func (r *RegistraSrvController) RegisterClass(key string, clazz any) {
	// Não verifico se existe, pois não é um add, ele simplesmente substitui o valor
	r.RegisteredClasses[key] = reflect.TypeOf(clazz)
}
func (r *RegistraSrvController) PrintRegisterClass() {
	for key := range r.RegisteredClasses {
		lib.NewLog().Info("nenhum", 0, 0, "Classe:", key, "Registrada com Sucesso")
	}
}

func (r *RegistraSrvController) FindSrvClassByKeyAndNewAsObject(key string) (IControllerSrv, error) {

	instancia, err := CriarStruct(key, r.RegisteredClasses)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return instancia.(IControllerSrv), nil
}

func GetRegistraSrvInstance() *RegistraSrvController {
	if registraSrvGlobal == nil {
		registraSrvGlobal = &RegistraSrvController{
			RegisteredClasses: make(RegisterClasses),
		}
	}
	return registraSrvGlobal
}

func RegistrarSrv(key string, clazz any) {
	GetRegistraSrvInstance().RegisterClass(key, clazz)
}
