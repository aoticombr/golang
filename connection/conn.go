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
func (co *Conn) Commit(tx *sql.Tx) {
	tx.Commit()
}
func (co *Conn) Rollback(tx *sql.Tx) {
	tx.Rollback()
}
func (co *Conn) Exec(tx *sql.Tx, sql string, arg ...any) (sql.Result, error) {
	return tx.Exec(sql, arg)
}

func (co *Conn) SetConfig(cf *ConfigOra) *Conn {
	co.config = cf
	return co
}
func (co *Conn) GetDB() *sql.DB {
	return co.db
}
func (co *Conn) Disconnect() {
	co.db.Close()
}

func GetConn(d Drive) (*Conn, error) {
	conn := &Conn{}
	conn.SetConfig(GetConfigOra(d).Load())
	db, err := sql.Open("oracle", conn.config.GetUrl())
	if err != nil {
		return nil, err
	}
	conn.db = db
	return conn, nil
}
