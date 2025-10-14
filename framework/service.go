package core

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aoticombr/golang/config"
	"github.com/aoticombr/golang/dbconndataset"
	"github.com/aoticombr/golang/framework/api"
	"github.com/aoticombr/golang/framework/bot"

	"github.com/aoticombr/golang/framework/srv"

	lib "github.com/aoticombr/golang/framework/lib"
	"github.com/joho/godotenv"
	"github.com/kardianos/service"
)

var onceApp sync.Once
var instanceApp *App
var logOS service.Logger
var CoreApi *api.CoreApi
var CoreSrv *srv.CoreSrv
var CoreBot *bot.CoreBot

type App struct {
	Name    string
	Config  *config.Config
	options []OptionsApp
}

func NewApp() *App {
	onceApp.Do(
		func() {
			instanceApp = &App{
				Name: lib.GetApplicationName(),
			}
		})

	return instanceApp
}

func NewAppDev(name string) *App {
	onceApp.Do(
		func() {
			instanceApp = &App{
				Name: name,
			}
		})

	return instanceApp
}

func (app *App) Execute(options ...OptionsApp) {
	app.options = options

	/*Se estiver rodando em sistema windows
	valida se o retorno da pasta atual contem a palavra system32
	pois se tiver é quase certeza que esta rodando como servico
	e caso tiver ele alterar a pasta do os como a pasta onde
	esta o exe
	Motivo: se nao fizer isso ele nao acha os arquivos locais*/
	if runtime.GOOS == "windows" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal("Erro ao obter o diretório atual:" + err.Error())
		}
		if strings.Contains(strings.ToUpper(dir), "SYSTEM32") {
			exePath, err := os.Executable()
			if err != nil {
				log.Fatal(err)
			}

			appDir := filepath.Dir(exePath)
			err = os.Chdir(appDir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	svcConfig := &service.Config{
		Name:        app.Name,
		DisplayName: lib.GetApplicationName() + " Go lang application AOTI",
		Description: lib.GetApplicationName() + " Go lang application AOTI",
	}

	// Instancia o serviço de execução de uma SRV ou API
	s, err := service.New(app, svcConfig)

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		emTeste := false
		for i := 0; i < len(os.Args); i++ {
			if os.Args[i] != "-test.run" {
				emTeste = true
			}
		}
		// Evita problemas no modo teste do VSCODE
		if !emTeste {
			err = service.Control(s, os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	logOS, err = s.Logger(nil)

	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()

	if err != nil {
		logOS.Error(err)
	}
}

func (app *App) Start(s service.Service) error {
	go app.Run()
	return nil
}

func (app *App) Stop(s service.Service) error {
	//	for _, core := range app.CoresSrv {
	//		core.Parar()
	//	}
	// Stop should not block. Return with a few seconds.
	<-time.After(time.Second * 2)
	return nil
}
func (app *App) Run() error {
	lib.NewLog().TypePrint = lib.LG_Silent
	lib.NewLog().AppName = app.Name
	// Le o json de configuração
	lib.NewLog().Screen(app.Name, "Lendo Json...")
	app.Config = config.NewConfig()

	// Defini os niveis de log da aplicação
	if app.Config.Log.Screen {
		lib.NewLog().TypePrint = lib.LG_PrintLn
	}

	lib.NewLog().Screen(app.Name, "Setando nivel de log...")
	lib.NewLog().Screen(app.Name, "Nivel:", app.Config.Log.Nivel)
	switch app.Config.Log.Nivel {
	case 0:
		lib.NewLog().LogNivel = lib.N_NENHUM
	case 1:
		lib.NewLog().LogNivel = lib.N_DEBUG
	case 2:
		lib.NewLog().LogNivel = lib.N_INFO
	case 3:
		lib.NewLog().LogNivel = lib.N_WARNING
	case 4:
		lib.NewLog().LogNivel = lib.N_ERROR
	case 5:
		lib.NewLog().LogNivel = lib.N_CRITICAL
	}
	for _, option := range app.options {
		go option(app)
	}

	return nil
}

func (app *App) dbconn(dbname string) (*config.Database, *dbconndataset.ConnDataSet) {

	ok, dbconf := app.Config.GetDB(dbname)
	if !ok {
		log.Panicf("banco de dados (%s) não localizado no config.json: ", dbname)
	}
	var (
		conn *dbconndataset.ConnDataSet
		err  error
	)

	conn, err = dbconndataset.NewConn(*dbconf)

	if err != nil {
		log.Panicf("erro ao conectar no banco de dados (%s): %s", dbname, err.Error())
	}

	conn.PoolSize = dbconf.PoolSize
	conn.MaxOpenConns = dbconf.MaxConn

	return dbconf, conn
}

func (app *App) StopService() {}
