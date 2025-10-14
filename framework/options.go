package core

import (
	"fmt"
	"time"

	"github.com/aoticombr/golang/framework/api"
	"github.com/aoticombr/golang/framework/bot"
	lib "github.com/aoticombr/golang/framework/lib"
	"github.com/aoticombr/golang/framework/srv"
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
				go CoreApi.Start()
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
			go CoreBot.Executar()

		}
	}
}
func AddSrv() OptionsApp {
	return func(app *App) {
		srv.GetRegistraSrvInstance().PrintRegisterClass()
		ok, srvConfig := app.Config.GetService(app.Name)
		if ok == false {
			lib.NewLog().Info(app.Name, "Servidor de Serviço não sera iniciado, sem configuração!!!")
		} else {
			lib.NewLog().Info(app.Name, "Iniciando CoresSrv...")

			CoreSrv = srv.NewCoreSrv(app.Config.Certs, app.Config.Dbs, srvConfig)

			for {
				lib.NewLog().Debug("[Api]", "Core.Srv.Executar()")
				go CoreSrv.Executar()

				time.Sleep(1 * time.Minute)
			}
		}
	}
}
