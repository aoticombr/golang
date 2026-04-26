package core

import (
	"fmt"

	"github.com/aoticombr/golang/framework/api"
	"github.com/aoticombr/golang/framework/bot"
	"github.com/aoticombr/golang/framework/monitor"
	"github.com/aoticombr/golang/framework/srv"
	lib "github.com/aoticombr/golang/lib"
)

type OptionsApp func(*App)

func AddApi() OptionsApp {
	return func(app *App) {
		ok, apiConfig := app.Config.GetApi(app.Name)

		if ok == false {
			lib.NewLog().Info(app.Name, "Servidor de Api não sera iniciado, sem configuração!!!")
		} else {
			if ok {
				lib.NewLog().Info(app.Name, "Iniciando CoresAPI...")

				CoreApi = api.NewCoreApi(app.Config.Certs, app.Config.Dbs, apiConfig, app.Config.Parans)

				processo := CoreApi.Processo

				for i, dbname := range apiConfig.Dbs {
					_, conn := app.dbconn(dbname)
					if i == 0 {
						processo.Carregar(conn)
					}
					fmt.Printf("add tenant")
					processo.Tenants.Add(dbname, conn)
				}

				lib.NewLog().Debug("[Api]", "Core.API.Start()")
				go func() {
					if err := CoreApi.Start(); err != nil {
						lib.NewLog().Error(app.Name, "CoreApi.Start finalizou com erro:", err.Error())
					}
				}()
			}
		}
	}
}
func AddBot() OptionsApp {
	return func(app *App) {
		bot.GetRegistraBotInstance().PrintRegisterClass()
		ok, botConfig := app.Config.GetBot(app.Name)
		if ok == false {
			lib.NewLog().Info(app.Name, "Servidor de Bots não sera iniciado, sem configuração!!!")
		} else {
			lib.NewLog().Info(app.Name, "Iniciando CoresBot...")

			CoreBot = bot.NewCoreBot(app.Config.Certs, app.Config.Dbs, botConfig)

			lib.NewLog().Debug("[Api]", "Core.Srv.Executar()")
			go func() {
				if err := CoreBot.Executar(app.Ctx); err != nil {
					lib.NewLog().Error(app.Name, "Erro ao iniciar CoreBot:", err.Error())
				}
			}()

		}
	}
}
func AddMonitor() OptionsApp {
	return func(app *App) {
		if app.Config.Monitor == nil || !app.Config.Monitor.Ativo {
			lib.NewLog().Info(app.Name, "Monitor não sera iniciado, sem configuração ou desativado!!!")
			return
		}
		lib.NewLog().Info(app.Name, "Iniciando Monitor...")
		m, err := monitor.New(app.Config, app.Name)
		if err != nil {
			lib.NewLog().Error(app.Name, "Falha ao iniciar Monitor:", err.Error())
			return
		}
		CoreMonitor = m
		go func() {
			if err := m.Start(app.Ctx); err != nil {
				lib.NewLog().Error(app.Name, "Monitor finalizou com erro:", err.Error())
			}
		}()
	}
}
func AddSrv() OptionsApp {
	return func(app *App) {
		srv.GetRegistraSrvInstance().PrintRegisterClass()
		ok, srvConfig := app.Config.GetService(app.Name)
		if !ok {
			lib.NewLog().Info(app.Name, "Servidor de Serviço não sera iniciado, sem configuração!!!")
			return
		}
		lib.NewLog().Info(app.Name, "Iniciando CoresSrv...")

		CoreSrv = srv.NewCoreSrv(app.Config.Certs, app.Config.Dbs, srvConfig)

		lib.NewLog().Debug("[Api]", "Core.Srv.Executar()")
		if err := CoreSrv.Executar(app.Ctx); err != nil {
			lib.NewLog().Error(app.Name, "Erro ao iniciar CoreSrv:", err.Error())
		}
	}
}
