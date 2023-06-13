package connection

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Drive int

var (
	ORA Drive = 0
	PG  Drive = 1
)
var cfgglobal *ConfigOra

type ConfigOra struct {
	Drive    Drive
	Host     string
	User     string
	Pass     string
	Port     int
	Database string
	Schema   string
	Sid      string
}

func (cf *ConfigOra) Load() *ConfigOra {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	switch cf.Drive {
	case ORA:
		cf.Host = GetEnvString("ora_host")
		cf.User = GetEnvString("ora_user")
		cf.Port = GetEnvInt("ora_port")
		cf.Pass = GetEnvString("ora_pass")
		cf.Schema = GetEnvString("ora_schema")
		cf.Sid = GetEnvString("ora_sid")
	case PG:
		cf.Host = GetEnvString("pg_host")
		cf.User = GetEnvString("pg_user")
		cf.Port = GetEnvInt("pg_port")
		cf.Pass = GetEnvString("pg_pass")
		cf.Database = GetEnvString("pg_database")
	}
	return cf
}

func (cf *ConfigOra) getUrlOra() string {
	url := fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		cf.User, cf.Pass, cf.Host, cf.Port, cf.Sid)
	return url
}
func (cf *ConfigOra) getUrlPG() string {
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cf.Host, cf.Port, cf.User, cf.Pass, cf.Database)
	return url
}
func (cf *ConfigOra) GetUrl() string {
	switch cf.Drive {
	case ORA:
		return cf.getUrlOra()
	case PG:
		return cf.getUrlPG()
	}
	return ""
}

func GetConfigOra(d Drive) *ConfigOra {
	if cfgglobal == nil {
		cfgglobal = &ConfigOra{}
		cfgglobal.Drive = d
	}
	return cfgglobal
}
