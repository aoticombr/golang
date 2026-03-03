package dbconndataset

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aoticombr/golang/config"
	dbconnect "github.com/aoticombr/golang/dbconnbase"
	"github.com/aoticombr/golang/dbdataset"
)

type ConnDataSet dbconnect.Conn

func (cds *ConnDataSet) Conn() *dbconnect.Conn {
	return (*dbconnect.Conn)(cds)
}
func (cds *ConnDataSet) Close() error {
	return cds.Conn().DB.Close()
}
func (cds *ConnDataSet) Ping() error {
	return cds.Conn().DB.Ping()
}

func (cds *ConnDataSet) NewDataSet() *dbdataset.DataSet {
	ds := dbdataset.NewDataSet(cds.Conn())
	return ds
}
func (cds *ConnDataSet) NewDataSetName(name string) *dbdataset.DataSet {
	ds := dbdataset.NewDataSet(cds.Conn())
	ds.Name = name
	return ds
}

// Métodos delegados do dbconnect.Conn
func (cds *ConnDataSet) Open() error {
	return (*dbconnect.Conn)(cds).Open()
}

func (cds *ConnDataSet) SetSizePool(n int) {
	(*dbconnect.Conn)(cds).SetSizePool(n)
}

func (cds *ConnDataSet) SetPoolLifeTime(d time.Duration) {
	(*dbconnect.Conn)(cds).SetPoolLifeTime(d)
}

func (cds *ConnDataSet) SetMaxOpenConns(n int) {
	(*dbconnect.Conn)(cds).SetMaxOpenConns(n)
}

func (cds *ConnDataSet) SetConnLifeTime(d time.Duration) {
	(*dbconnect.Conn)(cds).SetConnLifeTime(d)
}

func (cds *ConnDataSet) CreateContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return (*dbconnect.Conn)(cds).CreateContext(ctx)
}

func (cds *ConnDataSet) StartTransaction() (*dbconnect.Transaction, error) {
	return (*dbconnect.Conn)(cds).StartTransaction()
}

func (cds *ConnDataSet) StartTransactionContext(ctx context.Context) (*dbconnect.Transaction, error) {
	return (*dbconnect.Conn)(cds).StartTransactionContext(ctx)
}

func (cds *ConnDataSet) Exec(sql string, arg ...any) (sql.Result, error) {
	return (*dbconnect.Conn)(cds).Exec(sql, arg...)
}

func NewConn(db config.Database) (*ConnDataSet, error) {
	fmt.Println("Criando Connection")
	conn := &ConnDataSet{
		Dialect:  dbconnect.DialectLowFromString(db.Db),
		DSN:      db.GetDsn(),
		PoolSize: db.PoolSize,
	}

	fmt.Println("abrindo conexão")
	err := conn.Open() // Agora pode usar diretamente

	if err != nil {
		return nil, err
	}

	return conn, nil
}
