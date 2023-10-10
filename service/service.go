package service

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	service "github.com/kardianos/service"
)

var servicoglobal *servico
var logger service.Logger

type servico struct {
	Name        string
	Description string
	Logger      service.Logger
}

func (serv *servico) Exec(prog service.Interface) {
	/*Se estiver rodando em sistema windows
	valida se o retorno da pasta atual contem a palavra system32
	pois se tiver é quase certeza que esta rodando como servico
	e caso tiver ele alterar a pasta do os como a pasta onde
	esta o exe
	Motivo: se nao fizer isso ele nao acha o .env*/
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
	svcConfig := &service.Config{
		Name:        serv.Name,
		DisplayName: serv.Description,
		Description: serv.Description,
	}
	s, err := service.New(prog, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	serv.Logger = logger
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func NewServico() *servico {
	if servicoglobal == nil {
		servicoglobal = &servico{}
	}
	return servicoglobal
}
