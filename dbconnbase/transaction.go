package dbconnbase

import (
	"context"
	"database/sql"
)

type Transaction struct {
	Conn *Conn
	Tx   *sql.Tx
	Ctx  context.Context
}

func NewTransaction(conn *Conn) (*Transaction, error) {
	tx, err := conn.DB.Begin()

	if err != nil {
		return nil, err
	}

	transac := &Transaction{
		Conn: conn,
		Tx:   tx,
	}

	return transac, nil
}

func NewTransactionCtx(conn *Conn, ctx context.Context) (*Transaction, error) {
	opts := &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	}

	tx, err := conn.DB.BeginTx(ctx, opts)

	if err != nil {
		return nil, err
	}

	transac := &Transaction{
		Conn: conn,
		Tx:   tx,
		Ctx:  ctx,
	}

	return transac, nil
}

func (t *Transaction) Commit() error {
	return t.Tx.Commit()
}

func (t *Transaction) Rollback() error {
	return t.Tx.Rollback()
}

// func (t *Transaction) NewDataSet() *DataSet {
// 	ds := NewDataSetTx(t)
// 	ds.AddContext(t.Ctx)
// 	return ds
// }
