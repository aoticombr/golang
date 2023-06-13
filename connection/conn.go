package connection

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/sijms/go-ora/v2"
)

type Conn struct {
	config *ConfigOra
	db     *sql.DB
}

func (co *Conn) StartTransaction() (*sql.Tx, error) {
	return co.db.Begin()
}
func (co *Conn) Commit() error {
	err := co.tx.Commit()
	co.tx = nil
	return err
}
func (co *Conn) Rollback() error {
	err := co.tx.Rollback()
	co.tx = nil
	return err
}
func (co *Conn) Exec(sql string, arg ...any) (sql.Result, error) {
	return co.tx.Exec(sql, arg)
}

func (co *Conn) SetConfig(cf *ConfigOra) *Conn {
	co.config = cf
	return co
}

func (co *Conn) Disconnect() {
	co.Db.Close()
}

func GetConn(d Drive) (*Conn, error) {
	conn := &Conn{}
	conn.SetConfig(GetConfigOra(d).Load())
	db, err := sql.Open("oracle", conn.config.GetUrl())
	if err != nil {
		return nil, err
	}
	conn.Db = db
	return conn, nil
}
