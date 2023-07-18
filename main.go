package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/aoticombr/golang/Logger"
	cp "github.com/aoticombr/golang/component"
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
	for i := 0; i < 2; i++ {
		_, err := q2.ExecDirect()
		if err != nil {
			fmt.Println("i:", i, "erro:", err)
		}
		//time.Sleep(1 * time.Second)
		fmt.Println("CDMARCAS GetData int64:", q2.ParamByName("CDMARCAS").AsInt64())
		fmt.Println("CDMARCAS GetData string:", q2.ParamByName("CDMARCAS").AsString())
	}
	cl1, cl2 := cp.ConvertToInsertStatement(q2.Params)
	fmt.Println(cl1)
	fmt.Println(cl2)

	q3 := ds.GetDataSet(conn)
	q3.Sql.
		Add("SELECT CDMARCAS, MARCAS,ATIVO FROM NEXUS.MARCAS WHERE ATIVO = :ATIVO and Rownum <= 5")

	//q3.SetInputParam("MARCAS", "XXX")
	q3.SetInputParam("ATIVO", "S")
	q3.Open()
	for !q3.Eof() {
		fmt.Println("CDMARCASCDMARCAS GetData int64:", q3.FieldByName("CDMARCAS").AsInt64())
		fmt.Println("MARCAS GetData string:", q3.FieldByName("MARCAS").AsString())
		q3.Next()
	}

}
