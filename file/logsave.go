package file

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logGlobal *Log

type Log struct {
	save       bool
	viewscreen bool
}

func (l Log) SaveLog(texto ...string) {
	nome_arq := time.Now().Format("2006-01-02") + ".log"
	line := fmt.Sprintf("%s \n", texto)
	arquivo, err := os.OpenFile(nome_arq, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer arquivo.Close()

	if _, err := arquivo.WriteString(line); err != nil {
		log.Fatal(err.Error())
	}
}
func (l Log) SendMsg(texto string) {
	dataHoraAtual := time.Now().Format("2006-01-02 15:04:05")
	logString := fmt.Sprintf("[%s] %s\n", dataHoraAtual, texto)
	if l.viewscreen {
		fmt.Printf(logString)
	}
	if l.save {
		l.SaveLog(logString)
	}
}
func (l Log) SendErro(texto string, err error) {
	dataHoraAtual := time.Now().Format("2006-01-02 15:04:05")
	logString := fmt.Sprintf("[%s] Erro(%s):>> %s\n", dataHoraAtual, texto, err.Error())
	if l.viewscreen {
		fmt.Printf(logString)
	}
	if l.save {
		l.SaveLog(logString)
	}
}
func GetLog() *Log {
	if logGlobal == nil {
		logGlobal = &Log{}
	}
	return logGlobal
}
