package ora

import (
	"database/sql"

	_ "github.com/sijms/go-ora/v2"
)

type ConnOra struct {
	config *ConfigOra
	db     *sql.DB
}

func (co *ConnOra) SetConfig(cf *ConfigOra) *ConnOra {
	co.config = cf
	return co
}
func (co *ConnOra) GetDB() *sql.DB {
	return co.db
}

func GetConnOra() *ConnOra {
	conn := &ConnOra{}
	conn.SetConfig(GetConfigOra().Load())
	db, err := sql.Open("oracle", conn.config.GetUrl())
	if err != nil {
		panic("Error ao abrir conexao: " + err.Error())
	}
	conn.db = db
	return conn
}
