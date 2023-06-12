package main

import (
	"fmt"

	ora "github.com/aoticombr/golang/connection"
	ds "github.com/aoticombr/golang/dataset"
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
	conn, _ := ora.GetConn(ora.ORA)
	defer conn.Disconnect()
	q := ds.GetDataSet(conn)
	q.Sql.Add("SELECT * FROM CONTASPAGAR where rownum <= 10")
	q.Open()
	q.First()
	fmt.Println("q.Eof():", q.Eof())
	for !q.Eof() {
		fmt.Println(
			q.FieldByName("valor").AsFloat(),
			q.FieldByName("valor").AsString(),
		)
		q.Next()
	}

	// executablePath, err := os.Executable()
	// if err != nil {
	// 	// Lidar com o erro, se necessário
	// }
	// appRoot := filepath.Dir(executablePath)
	// logDir := appRoot //
	// fmt.Println(logDir)
	// logger, _ := log.NewLogger("INFO", os.Stdout, "[DEVRAIZ]", logDir)
	// logger.Info("ler o arquivo")

}
