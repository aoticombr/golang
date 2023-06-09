package dataset

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
)

type DataSet struct {
	DB               *sql.DB
	Sql              []string
	Rows             []map[string]Field
	Param            map[string]Parameter
	Eof              bool
	Index            int
	Recno            int
	Count            int
	DetailFields     string
	MasterSouce      *DataSet
	MasterFields     string
	MasterDetailList map[string]MasterDetails
	IndexFieldNames  string
}

func GetDataSet(db *sql.DB) *DataSet {
	var dataSet DataSet

	dataSet.DB = db

	dataSet.Index = 0
	dataSet.Recno = 0
	dataSet.Count = 0
	dataSet.Eof = true
	dataSet.Param = make(map[string]Parameter)

	return &dataSet
}

func (ds *DataSet) Open() error {
	ds.Rows = nil
	ds.Index = 0
	ds.Recno = 0
	ds.Count = 0
	ds.Eof = true

	vsql := ds.GetSql()

	var param []any
	for _, prm := range ds.Param {
		param = append(param, prm.value)
	}

	rows, err := ds.DB.Query(vsql, param...)

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

	ds.Scan(rows)

	ds.Count = len(ds.Rows)

	if ds.Count > 0 {
		ds.Recno = 1
		ds.Eof = ds.Count == 1
	}
	return nil
}

func (ds *DataSet) StartTransaction() (*sql.Tx, error) {
	tx, err := ds.DB.Begin()
	return tx, err
}
func (ds *DataSet) Commit(tx *sql.Tx) {
	tx.Commit()
}
func (ds *DataSet) Rollback(tx *sql.Tx) {
	tx.Rollback()
}
func (ds *DataSet) ExecTransact(tx *sql.Tx) (sql.Result, error) {
	var param []any
	for _, prm := range ds.Param {
		param = append(param, prm.value)
	}
	res, err := tx.Exec(ds.GetSql(), param...)
	return res, err
}
func (ds *DataSet) ExecDirect() (sql.Result, error) {
	var param []any
	for _, prm := range ds.Param {
		param = append(param, prm.value)
	}
	result, err := ds.DB.Exec(ds.GetSql(), param...)
	return result, err
}

func (ds *DataSet) AddSql(sql string) *DataSet {
	ds.Sql = append(ds.Sql, sql)

	return ds
}
func (ds *DataSet) ClearSql() {
	ds.Sql = nil
}

func (ds *DataSet) GetSql() (sql string) {
	for i, s := range ds.Sql {
		if i != len(ds.Sql)-1 {
			sql = sql + s + " \n"
		} else {
			sql = sql + s
		}
	}

	if ds.MasterSouce != nil {
		var sqlWhereMasterDetail string
		mf := strings.Split(ds.MasterFields, ";")
		df := strings.Split(ds.DetailFields, ";")

		for i := 0; i < len(mf); i++ {
			aliasHash, _ := uuid.NewUUID()
			alias := strings.Replace(aliasHash.String(), "-", "", -1)
			if i == len(mf)-1 {
				sqlWhereMasterDetail = sqlWhereMasterDetail + df[i] + " = :" + alias
			} else {
				sqlWhereMasterDetail = sqlWhereMasterDetail + df[i] + " = :" + alias + " and "
			}

			ds.ParamByName(alias, ds.MasterSouce.FieldByName(mf[i]).Value)
		}

		if sqlWhereMasterDetail != "" {
			sql = "select * from (" + sql + ") where " + sqlWhereMasterDetail
		}
	}

	return sql
}

func (ds *DataSet) Scan(list *sql.Rows) {
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

		row := make(map[string]Field)

		for i, value := range columns {
			row[fields[i]] = Field{name: fields[i],
				caption:    fields[i],
				dataType:   columntypes[i],
				Value:      value,
				dataMask:   "",
				valueTrue:  "",
				valueFalse: "",
				visible:    true,
				order:      i + 1,
				index:      i,
			}
		}

		ds.Rows = append(ds.Rows, row)
	}
}

func (ds *DataSet) ParamByName(paramName string, paramValue any) *DataSet {

	ds.Param[paramName] = Parameter{value: paramValue}

	return ds
}

func (ds *DataSet) FieldByName(fieldName string) Field {
	field := strings.ToUpper(fieldName)
	return ds.Rows[ds.Index][field]
}

func (ds *DataSet) Locate(key string, value any) bool {

	ds.First()
	for ds.Eof == false {
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
	ds.Index = 0
	ds.Recno = 0
	if ds.Count > 0 {
		ds.Index = 0
		ds.Recno = 1
		ds.Eof = ds.Count == 0
	} else {
		ds.Eof = true
	}
}

func (ds *DataSet) Next() {
	if !ds.Eof {
		if ds.Recno < ds.Count {
			ds.Eof = ds.Count == ds.Recno
			ds.Index++
			ds.Recno++
		} else {
			ds.Eof = true
		}
	}
}

func (ds *DataSet) IsEmpty() bool {
	return ds.Count == 0
}

func (ds *DataSet) IsNotEmpty() bool {
	return ds.Count > 0
}
