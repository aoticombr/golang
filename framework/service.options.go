package core

import (
	"fmt"
	"strings"

	"github.com/kardianos/service"
)

// ServiceOptions personaliza o registro do aplicativo como serviço do SO.
// Campos vazios mantêm o comportamento padrão do framework.
type ServiceOptions struct {
	// DisplayName é o nome exibido no gerenciador de serviços.
	DisplayName string
	// Description é a descrição do serviço.
	Description string
	// Dependencias são outros serviços que precisam subir antes
	// (formato do kardianos: "After=network.target" no Linux, nomes no Windows).
	Dependencias []string
	// Usuario é a conta que executa o serviço (vazio = padrão do SO: LocalSystem/root).
	Usuario string
	// Argumentos passados ao executável quando iniciado pelo gerenciador de serviços.
	Argumentos []string
	// Restart (Linux systemd): "always", "on-failure"... vazio = padrão do kardianos.
	Restart string
	// InicioManual (Windows): true = início Manual; false (padrão) = Automático.
	InicioManual bool
}

// serviceOptionsPendentes guarda a configuração informada por AddServiceConfig.
// AddServiceConfig é avaliada como argumento de Execute, antes do service.New,
// por isso o valor é guardado aqui e consumido por aplicarServiceOptions.
var serviceOptionsPendentes *ServiceOptions

// AddServiceConfig personaliza o serviço do SO (nome exibido, descrição,
// dependências, conta, restart). Opcional: sem ela o framework mantém
// DisplayName/Description "<exe> Go lang application AOTI".
//
//	core.NewAppDev("meuapp").Execute(
//	    core.AddServiceConfig(core.ServiceOptions{DisplayName: "Meu App", Restart: "always"}),
//	    core.AddApi(),
//	)
func AddServiceConfig(opts ServiceOptions) OptionsApp {
	o := opts
	serviceOptionsPendentes = &o
	return func(app *App) {} // nada a fazer em tempo de execução
}

// ServiceConfig é a forma encadeada de AddServiceConfig:
//
//	core.NewAppDev("meuapp").ServiceConfig(core.ServiceOptions{...}).Execute(core.AddApi())
func (app *App) ServiceConfig(opts ServiceOptions) *App {
	o := opts
	serviceOptionsPendentes = &o
	return app
}

func (app *App) aplicarServiceOptions(cfg *service.Config) {
	o := serviceOptionsPendentes
	if o == nil {
		return
	}
	if o.DisplayName != "" {
		cfg.DisplayName = o.DisplayName
	}
	if o.Description != "" {
		cfg.Description = o.Description
	}
	if len(o.Dependencias) > 0 {
		cfg.Dependencies = o.Dependencias
	}
	if o.Usuario != "" {
		cfg.UserName = o.Usuario
	}
	if len(o.Argumentos) > 0 {
		cfg.Arguments = o.Argumentos
	}
	if cfg.Option == nil {
		cfg.Option = service.KeyValue{}
	}
	if o.Restart != "" {
		cfg.Option["Restart"] = o.Restart
	}
	if o.InicioManual {
		cfg.Option["StartType"] = "manual"
	}
}

// verbosServico são os comandos aceitos como primeiro argumento.
var verbosServico = map[string]bool{
	"install": true, "uninstall": true, "start": true,
	"stop": true, "restart": true, "status": true,
}

// VerboServico informa se a linha de comando pede uma ação de controle do
// serviço. Execuções de teste (argumentos -test.*) nunca são tratadas como verbo.
func VerboServico(args []string) (string, bool) {
	for _, a := range args[1:] {
		if strings.HasPrefix(a, "-test.") {
			return "", false
		}
	}
	if len(args) < 2 {
		return "", false
	}
	v := strings.ToLower(strings.TrimSpace(args[1]))
	if !verbosServico[v] {
		return "", false
	}
	return v, true
}

func controlarServico(s service.Service, verbo string) error {
	if verbo == "status" {
		st, err := s.Status()
		if err != nil {
			return err
		}
		fmt.Println(statusTexto(st))
		return nil
	}
	if err := service.Control(s, verbo); err != nil {
		return err
	}
	fmt.Println("serviço:", verbo, "ok")
	return nil
}

func statusTexto(st service.Status) string {
	switch st {
	case service.StatusRunning:
		return "executando"
	case service.StatusStopped:
		return "parado"
	default:
		return "desconhecido"
	}
}
