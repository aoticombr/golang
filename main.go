package main

import (
	"fmt"
	"time"

	ora "github.com/aoticombr/golang/connection"
	ds "github.com/aoticombr/golang/dataset"
)

type ContasPagar struct {
	Cdcontaspagar int64
	Historico     string
	Dtaconta      time.Time
	Valor         float64
}

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
	q.Sql.Add("SELECT cdcontaspagar, historico, dtaconta, valor FROM CONTASPAGAR where rownum <= 10")
	q.Open()
	q.First()
	// fmt.Println("q.Eof():", q.Eof())
	// for !q.Eof() {
	// 	fmt.Println(
	// 		q.FieldByName("cdcontaspagar").AsInt64(),
	// 		q.FieldByName("historico").AsString(),
	// 		q.FieldByName("dtaconta").AsDateTime(),
	// 		q.FieldByName("valor").AsFloat64(),
	// 	)
	// 	q.Next()
	// }
	//var contas []ContasPagar
	results, err := q.RowInStruck(ContasPagar{})
	if err != nil {
		fmt.Println("Erro ao executar a consulta:", err)
		return
	}
	fmt.Println("results:", results)
	for _, result := range results {
		contasPagar, ok := result.(ContasPagar) // Faz a conversão para o tipo correto

		if !ok {
			fmt.Println("Erro ao converter resultado para ContasPagar")
			continue
		}

		fmt.Println("Cdcontaspagar:", contasPagar.Cdcontaspagar)
		fmt.Println("Historico:", contasPagar.Historico)
		fmt.Println("Dtaconta:", contasPagar.Dtaconta)
		fmt.Println("Valor:", contasPagar.Valor)

		fmt.Println("----------------------")
	}
	// for _, result := range results {
	// 	conta := result.(*ContasPagar)
	// 	fmt.Println(*conta)
	// }

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
