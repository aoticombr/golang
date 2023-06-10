package main

import (
	"fmt"

	ora "github.com/aoticombr/go/connection"
	ds "github.com/aoticombr/go/dataset"
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
	//defer conn.FreeAndNil()
	q := ds.GetDataSet(conn)
	q.Sql.Add("SELECT * FROM CONTASPAGAR where rownum <= 10")
	q.Open()
	q.First()
	for !q.Eof() {
		fmt.Println(
			//	q.FieldByName("cdcontaspagar").AsInt64(),
			q.FieldByName("valor").AsFloat(),
			//q.FieldByName("juros").AsString()
		)
		q.Next()
	}

}
