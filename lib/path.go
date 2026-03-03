package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/*
GetApplicationName:
nome da aplicação sem extensão
*/
func GetApplicationName() string {
	name := GetApplicationNameExe()
	//em caso de windows remove o exe do nome
	if runtime.GOOS == "windows" {
		name = strings.TrimSuffix(name, ".exe")
	}
	return name
}

/*
GetPathApplication:
Pasta da aplicação
*/
func GetPathApplication() string {

	Path, _ := os.Getwd()

	if runtime.GOOS == "windows" {
		if strings.Contains(Path, "system32") {
			exePath, err := os.Executable()

			if err != nil {
				fmt.Println("Erro ao obter o caminho do executável:", err)
			}

			Path = filepath.Dir(exePath)
		}
	}

	return Path
}

/*
GetPathApplicationFile:
Aplica pasta local da aplicação ao nome do arquivo informado
*/
func GetPathApplicationFile(fileName string) string {
	var file string
	file = filepath.Join(GetPathApplication(), fileName)
	return file
}
