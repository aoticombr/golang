package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/aoticombr/golang/Logger"
	ora "github.com/aoticombr/golang/connection"
	ds "github.com/aoticombr/golang/dataset"
)

type ContasPagar struct {
	Cdcontaspagar int64
	Historico     string
	Dtaconta      time.Time
	Valor         float32
}

type any struct {
	// Defina os campos necessários para o tipo `any`
}

type value interface {
	GetData() interface{} // Exemplo de um método na interface `value`
}

type component struct {
	Value value
	// Defina outros campos necessários para o componente
}

func main() {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	appRoot := filepath.Dir(executablePath)
	logDir := appRoot //
	fmt.Println(logDir)
	logger, _ := log.NewLogger("ERROR", os.Stdout, "[DEVRAIZ]", logDir)

	logger.Warning("Warning================")
	//logger.Fatal("Fatal================")
	conn, _ := ora.GetConn(ora.ORA)
	defer conn.Disconnect()

	q2 := ds.GetDataSet(conn)
	q2.Sql.
		Add("INSERT INTO NEXUS.MARCAS").
		Add("(MARCAS,ATIVO)").
		Add("VALUES").
		Add("(:MARCAS,:ATIVO)").
		Add("RETURNING CDMARCAS INTO :CDMARCAS")

	q2.SetInputParam("MARCAS", "XXX")
	q2.SetInputParam("ATIVO", "S")
	q2.SetOutputParam("CDMARCAS", int64(0))
	for i := 0; i < 288; i++ {
		_, err := q2.ExecDirect()
		if err != nil {
			fmt.Println("i:", i, "erro:", err)
		}
		//time.Sleep(1 * time.Second)
		fmt.Println("CDMARCAS GetData:", q2.ParamByName("CDMARCAS").GetData())
	}

}
