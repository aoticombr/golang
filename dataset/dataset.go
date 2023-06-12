package dataset

import (
	"database/sql"
	"strings"

	cp "github.com/aoticombr/golang/component"
	conn "github.com/aoticombr/golang/connection"
)

type DataSet struct {
	Connection *conn.Conn
	Sql        cp.Strings
	rows       cp.Rows
	param      cp.Params
	index      int
	Recno      int
}

func (ds *DataSet) Eof() bool {
	return ds.Count() == 0 || ds.Recno > ds.Count()
}
func (ds *DataSet) Count() int {
	return len(ds.rows)
}
func (ds *DataSet) GetParams() []any {
	var param []any
	for _, prm := range ds.param {
		param = append(param, prm.Value)
	}
	return param
}
func (ds *DataSet) Open() error {
	ds.rows = nil
	ds.index = 0
	ds.Recno = 0
	rows, err := ds.Connection.GetDB().Query(ds.Sql.Text(), ds.GetParams()...)

	if err != nil {
		return err
	}

	defer rows.Close()

	ds.scan(rows)

	ds.First()
	return nil
}

func (ds *DataSet) StartTransaction() (*sql.Tx, error) {
	return ds.Connection.GetDB().Begin()
}
func (ds *DataSet) Commit(tx *sql.Tx) {
	tx.Commit()
}
func (ds *DataSet) Rollback(tx *sql.Tx) {
	tx.Rollback()
}
func (ds *DataSet) ExecTransact(tx *sql.Tx) (sql.Result, error) {
	return tx.Exec(ds.Sql.Text(), ds.GetParams()...)
}
func (ds *DataSet) ExecDirect() (sql.Result, error) {
	return ds.Connection.GetDB().Exec(ds.Sql.Text(), ds.GetParams()...)
}
func (ds *DataSet) scan(list *sql.Rows) {
	columntypes, _ := list.ColumnTypes()
	fields, _ := list.Columns()
	for list.Next() {
		columns := make([]interface{}, len(fields))

		for i := range columns {
			columns[i] = &columns[i]
		}

		if err := list.Scan(columns...); err != nil {
			panic(err)
		}

		row := make(map[string]cp.Field)

		for i, value := range columns {
			f := cp.Field{
				Name:       fields[i],
				Caption:    fields[i],
				DataType:   columntypes[i],
				DataMask:   "",
				Value:      cp.Variant{Value: value},
				ValueTrue:  "",
				ValueFalse: "",
				Visible:    true,
				Order:      i + 1,
				Index:      i,
			}

			row[fields[i]] = f
		}

		ds.rows = append(ds.rows, row)
	}
}
func (ds *DataSet) ParamByName(paramName string, paramValue any) *DataSet {

	ds.param[paramName] = cp.Parameter{Value: paramValue}

	return ds
}
func (ds *DataSet) FieldByName(fieldName string) cp.Field {
	field := strings.ToUpper(fieldName)
	return ds.rows[ds.index][field]
}
func (ds *DataSet) Locate(key string, value any) bool {
	ds.First()
	for !ds.Eof() {

		switch v := value.(type) {
		case string:
			if ds.FieldByName(key).AsValue() == v {
				return true
			}
		default:
			if ds.FieldByName(key).AsValue() == v {
				return true
			}
		}
		ds.Next()
	}
	return false
}
func (ds *DataSet) First() {
	ds.index = 0
	ds.Recno = 0
	if ds.Count() > 0 {
		ds.Recno = 1
	}
}
func (ds *DataSet) Next() {
	if !ds.Eof() {
		ds.index++
		ds.Recno++
	}
}
func (ds *DataSet) IsEmpty() bool {
	return ds.Count() == 0
}
func (ds *DataSet) IsNotEmpty() bool {
	return ds.Count() > 0
}
func GetDataSet(pconn *conn.Conn) *DataSet {
	ds := &DataSet{
		Connection: pconn,
		index:      0,
		Recno:      0,
		param:      make(map[string]cp.Parameter),
	}
	return ds
}
