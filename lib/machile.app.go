package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
)

/*
GetApplicationName:
nome da aplicação sem extensão
*/
func GetComputerName() string {
	name, err := os.Hostname()
	if err != nil {
		fmt.Println("Erro ao obter o nome da maquina:", err)
	}
	return name
}

/*
GetApplicationNameExe:
nome da aplicação com extensão
*/
func GetApplicationNameExe() string {
	execPath, _ := os.Executable()
	_, execName := filepath.Split(execPath)
	return execName
}

func GetModuleName() string {
	bi, ok := debug.ReadBuildInfo()

	if !ok {
		log.Printf("Falha ao extrair o path")
		return ""
	}

	return bi.Path
}
