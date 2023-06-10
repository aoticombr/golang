package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/aoticombr/go/logger"
)

func main() {
	//logLevel := flag.String("log", "ERROR", "Logging level")
	//flag.Parse()
	//logger, _ := log.NewLogger(*logLevel, os.Stdout, "[DEVRAIZ]")
	// save.GetLog().SaveLog("aaaa", "aaaa", "aaaa", "aaaa", "aaaa", "aaaa", "aaaa", "aaaa")
	// save.GetLog().SaveLog("bbbb")
	// save.GetLog().SaveLog("ccc")
	// logger.Info("Download", "Download", "Download", "Download", "Download", "Download", "Download")
	// logger.Info("Descompactar o arquivo")
	// logger.Info("ler o arquivo")
	// logger.Fatal("erro ao ler o arquivo")
	// logger.Debug("Debug================")
	// logger.Warning("Warning================")
	// logger.Fatal("Fatal================")
	
	executablePath, err := os.Executable()
	if err != nil {
		// Lidar com o erro, se necessário
	}
	appRoot := filepath.Dir(executablePath)
	logDir := appRoot //
	fmt.Println(logDir)
	logger, _ := log.NewLogger("INFO", os.Stdout, "[DEVRAIZ]", logDir)
	logger.Info("ler o arquivo")

}
