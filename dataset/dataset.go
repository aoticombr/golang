package dataset

import (
	"database/sql"
	"strings"

	cp "github.com/aoticombr/go/component"
	conn "github.com/aoticombr/go/connection"
)

type DataSet struct {
	Connection *conn.Conn
	Sql        cp.Strings
	rows       cp.Rows
	param      cp.Params
	//eof        bool
	index int
	recno int
	count int
}

func (ds *DataSet) Eof() bool {
	eof := true
	if ds.count == 0 {
		return true
	}
	eof = ds.recno > ds.count
	return eof
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
	ds.recno = 0
	ds.count = 0
	rows, err := ds.Connection.GetDB().Query(ds.Sql.Text(), ds.GetParams()...)

	if err != nil {
		return err
	}

	defer func(rows *sql.Rows) error {
		err := rows.Close()
		if err != nil {
			return err
		}
		return nil
	}(rows)

	ds.scan(rows)

	ds.count = len(ds.rows)

	ds.First()
	return nil
}

func (ds *DataSet) StartTransaction() (*sql.Tx, error) {
	tx, err := ds.Connection.GetDB().Begin()
	return tx, err
}
func (ds *DataSet) Commit(tx *sql.Tx) {
	tx.Commit()
}
func (ds *DataSet) Rollback(tx *sql.Tx) {
	tx.Rollback()
}
func (ds *DataSet) ExecTransact(tx *sql.Tx) (sql.Result, error) {
	res, err := tx.Exec(ds.Sql.Text(), ds.GetParams()...)
	return res, err
}
func (ds *DataSet) ExecDirect() (sql.Result, error) {
	result, err := ds.Connection.GetDB().Exec(ds.Sql.Text(), ds.GetParams()...)
	return result, err
}
func (ds *DataSet) scan(list *sql.Rows) {
	columntypes, _ := list.ColumnTypes()
	fields, _ := list.Columns()
	for list.Next() {
		columns := make([]interface{}, len(fields))

		for i := range columns {
			columns[i] = &columns[i]
		}

		err := list.Scan(columns...)

		if err != nil {
			panic(err)
		}

		row := make(map[string]cp.Field)

		for i, value := range columns {
			row[fields[i]] = cp.Field{
				Name:       fields[i],
				Caption:    fields[i],
				DataType:   columntypes[i],
				Value:      value,
				DataMask:   "",
				ValueTrue:  "",
				ValueFalse: "",
				Visible:    true,
				Order:      i + 1,
				Index:      i,
			}
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
	for ds.Eof() == false {
		switch value.(type) {
		case string:
			if ds.FieldByName(key).Value == value {
				return true
			}
		default:
			if ds.FieldByName(key).Value == value {
				return true
			}
		}

		ds.Next()
	}
	return false
}
func (ds *DataSet) First() {
	ds.index = 0
	ds.recno = 0
	if ds.count > 0 {
		ds.recno = 1
	}
}
func (ds *DataSet) Next() {
	if !ds.Eof() {
		//if ds.recno < ds.count {
		ds.index++
		ds.recno++
		//}
	}
}
func (ds *DataSet) IsEmpty() bool {
	return ds.count == 0
}
func (ds *DataSet) IsNotEmpty() bool {
	return ds.count > 0
}
func GetDataSet(pconn *conn.Conn) *DataSet {
	ds := &DataSet{
		Connection: pconn,
		index:      0,
		recno:      0,
		count:      0,
		param:      make(map[string]cp.Parameter),
	}
	return ds
}
