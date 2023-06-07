package ora

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var cfgglobal *ConfigOra

type ConfigOra struct {
	host   string
	user   string
	pass   string
	port   string
	schema string
	sid    string
}

func (cf *ConfigOra) SetHost(value string) *ConfigOra {
	cf.host = value
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
	cf.SetHost(os.Getenv("ora_host"))
	cf.SetUser(os.Getenv("ora_user"))
	cf.SetPort(os.Getenv("ora_port"))
	cf.SetPassCrypt(os.Getenv("ora_pass"))
	cf.SetSchema(os.Getenv("ora_schema"))
	cf.SetSid(os.Getenv("ora_sid"))
	return cf
}

func (cf ConfigOra) GetUrl() string {
	url := fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		cf.GetUser(),
		cf.GetPass(),
		cf.GetHost(),
		cf.GetPort(),
		cf.GetSid())
	return url
}

func GetConfigOra() *ConfigOra {
	if cfgglobal == nil {
		cfgglobal = &ConfigOra{}
	}
	return cfgglobal
}
