package bot

import (
	"fmt"
	"reflect"

	"github.com/aoticombr/golang/framework/lib"
)

type BotClass struct {
	Classe reflect.Type
	Tempo  int
}

type RegisterBotClasses map[string]BotClass

type RegistraBotController struct {
	RegisteredClasses RegisterBotClasses
}

var registraBotGlobal *RegistraBotController

func CriarStruct(nome string, mapaDeTipos map[string]BotClass) (interface{}, error) {
	// Verifica se o nome está presente no mapa
	tipo, ok := mapaDeTipos[nome]
	if !ok {
		return nil, fmt.Errorf("Struct com o nome %s não encontrada", nome)
	}

	// Cria uma instância da struct como um ponteiro
	instancia := reflect.New(tipo.Classe).Interface()

	return instancia, nil
}

func (r *RegistraBotController) RegisterClass(key string, clazz any, tempo_minuto int) {
	// Não verifico se existe, pois não é um add, ele simplesmente substitui o valor
	r.RegisteredClasses[key] = BotClass{
		Classe: reflect.TypeOf(clazz),
		Tempo:  tempo_minuto,
	}
}
func (r *RegistraBotController) PrintRegisterClass() {
	for key := range r.RegisteredClasses {
		lib.NewLog().Info("nenhum", 0, 0, "Classe:", key, "Registrada com Sucesso")
	}
}

func (r *RegistraBotController) FindBotClassByKeyAndNewAsObject(key string) (IControllerBot, error) {

	instancia, err := CriarStruct(key, r.RegisteredClasses)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return instancia.(IControllerBot), nil
}

func GetRegistraBotInstance() *RegistraBotController {
	if registraBotGlobal == nil {
		registraBotGlobal = &RegistraBotController{
			RegisteredClasses: make(RegisterBotClasses),
		}
	}
	return registraBotGlobal
}

func RegistrarBot(key string, clazz any, tempo_minuto int) {
	GetRegistraBotInstance().RegisterClass(key, clazz, tempo_minuto)
}
