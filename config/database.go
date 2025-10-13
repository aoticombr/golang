package config

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Database struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Schema   string `json:"schema"`
	Sid      string `json:"sid"`
	PoolSize int    `json:"poolsize"`
	MaxConn  int    `json:"maxconn"`
	Ativo    bool   `json:"ativo"`
	Trace    Trace  `json:"trace"`
	Db       string `json:"db"`
}

func (d Database) getUrlOra() string {
	/*
		CONNECTION TIMEOUT = RESPONSAVE PELO TEMPO DE ABERTURA DA QUERY
		SE PASSAR DESSE TEMPO ELE MANDA UM COMANDO PARA DERRUBAR
	*/
	var (
		minutos     int64
		minutos_str string
	)
	minutos = int64(time.Minute * 60)            //60 minutos
	minutos_str = strconv.FormatInt(minutos, 10) //convertendo para string
	url := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", d.User, d.Pass, d.Host, d.Port, d.Sid)
	url += "/?connection timeout=" + minutos_str + "&lob fetch=post"
	if d.Trace.Ativo {
		url += "&TRACE DIR=" + d.Trace.Path
	}
	return url
}

func (d Database) getUrlPg() string {
	// "postgres://postgres:manager@localhost:5432/nbs_status?sslmode=disable"
	senhaEscapada := url.QueryEscape(d.Pass)
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.User, senhaEscapada, d.Host, d.Port, d.Schema)
	return url
}
func (d Database) GetDsn() string {
	switch d.Db {
	case "ORA":
		return d.getUrlOra()
	case "PG":
		return d.getUrlPg()
	default:
		return ""
	}
}
