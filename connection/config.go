package ora

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Drive int

var (
	ORA Drive = 0
	PG  Drive = 1
)
var cfgglobal *ConfigOra

type ConfigOra struct {
	drive    Drive
	host     string
	user     string
	pass     string
	port     string
	database string
	schema   string
	sid      string
}

func (cf *ConfigOra) SetHost(value string) *ConfigOra {
	cf.host = value
	return cf
}
func (cf *ConfigOra) SetDatabase(value string) *ConfigOra {
	cf.database = value
	return cf
}
func (cf *ConfigOra) SetUser(value string) *ConfigOra {
	cf.user = value
	return cf
}
func (cf *ConfigOra) SetPassCrypt(value string) *ConfigOra {
	cf.pass = value
	return cf
}
func (cf *ConfigOra) SetPort(value string) *ConfigOra {
	cf.port = value
	return cf
}
func (cf *ConfigOra) SetSid(value string) *ConfigOra {
	cf.sid = value
	return cf
}
func (cf *ConfigOra) SetSchema(value string) *ConfigOra {
	cf.schema = value
	return cf
}
func (cf *ConfigOra) GetHost() string {
	return cf.host
}
func (cf *ConfigOra) GetUser() string {
	return cf.user
}
func (cf *ConfigOra) GetDatabase() string {
	return cf.database
}
func (cf *ConfigOra) GetPass() string {
	return cf.pass
}
func (cf *ConfigOra) GetSid() string {
	return cf.sid
}
func (cf *ConfigOra) GetPort() int {
	intValue, _ := strconv.Atoi(cf.port)
	return intValue
}
func (cf *ConfigOra) GetSchema() string {
	return cf.schema
}

func (cf *ConfigOra) Load() *ConfigOra {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	switch cf.drive {
	case ORA:
		cf.SetHost(os.Getenv("ora_host"))
		cf.SetUser(os.Getenv("ora_user"))
		cf.SetPort(os.Getenv("ora_port"))
		cf.SetPassCrypt(os.Getenv("ora_pass"))
		cf.SetSchema(os.Getenv("ora_schema"))
		cf.SetSid(os.Getenv("ora_sid"))
	case PG:
		cf.SetHost(os.Getenv("pg_host"))
		cf.SetUser(os.Getenv("pg_user"))
		cf.SetPort(os.Getenv("pg_port"))
		cf.SetPassCrypt(os.Getenv("pg_pass"))
		cf.SetDatabase(os.Getenv("pg_database"))
	}
	return cf
}

func (cf *ConfigOra) GetUrlOra() string {
	url := fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		cf.GetUser(), cf.GetPass(), cf.GetHost(), cf.GetPort(), cf.GetSid())
	return url
}
func (cf *ConfigOra) GetUrlPG() string {
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cf.GetHost(), cf.GetPort(), cf.GetUser(), cf.GetPass(), cf.GetDatabase())
	return url
}
func (cf *ConfigOra) GetUrl() string {
	switch cf.drive {
	case ORA:
		return cf.GetUrlOra()
	case PG:
		return cf.GetUrlPG()
	}
	return ""
}

func GetConfigOra(d Drive) *ConfigOra {
	if cfgglobal == nil {
		cfgglobal = &ConfigOra{}
		cfgglobal.drive = d
	}
	return cfgglobal
}
