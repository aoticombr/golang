package bot

import (
	"fmt"
	"reflect"
	"time"

	"github.com/aoticombr/golang/lib"
)

type BotClass struct {
	Classe reflect.Type
	// Tempo é o intervalo em minutos (RegistrarBot).
	Tempo int
	// Intervalo, quando > 0, tem precedência sobre Tempo (RegistrarBotIntervalo).
	Intervalo time.Duration
}

// Ciclo devolve o intervalo efetivo do bot. Zero/negativo vira 1 minuto
// (antes um Tempo 0 causava panic em time.NewTicker).
func (b BotClass) Ciclo() time.Duration {
	d := b.Intervalo
	if d <= 0 {
		d = time.Duration(b.Tempo) * time.Minute
	}
	if d <= 0 {
		d = time.Minute
	}
	return d
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

// RegistrarBotIntervalo registra um bot com intervalo livre (ex.: 30*time.Second).
// Mesmas regras de RegistrarBot: passe um VALOR (MeuBot{}) e implemente
// Executar() ou ExecutarCtx(ctx) no ponteiro.
func RegistrarBotIntervalo(key string, clazz any, intervalo time.Duration) {
	GetRegistraBotInstance().RegisteredClasses[key] = BotClass{
		Classe:    reflect.TypeOf(clazz),
		Intervalo: intervalo,
	}
}

// NovaInstancia cria a instância do bot sem exigir IControllerBot, permitindo
// bots que implementam apenas IControllerBotCtx.
func (r *RegistraBotController) NovaInstancia(key string) (any, error) {
	return CriarStruct(key, r.RegisteredClasses)
}
