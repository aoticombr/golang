package dataset

import (
	"database/sql"
	"fmt"
	"log"
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

func GetDataSetNew(db *sql.DB) *DataSet {
	var dataSet DataSet

	dataSet.DB = db

	dataSet.Index = 0
	dataSet.Recno = 0
	dataSet.Count = 0
	dataSet.Eof = true
	dataSet.Param = make(map[string]Parameter)

	return &dataSet
}

func (ds *DataSet) Open() {
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
		panic("Error to execute query. " + err.Error())
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("could not close the given rows %v\n", err)
		}
	}(rows)

	ds.Scan(rows)

	ds.Count = len(ds.Rows)

	if ds.Count > 0 {
		ds.Recno = 1
		ds.Eof = ds.Count == 1
	}
}

func (ds *DataSet) Exec() error {
	var param []any

	for _, prm := range ds.Param {
		param = append(param, prm.value)
	}

	result, err := ds.DB.Exec(ds.GetSql(), param...)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("error to execute query: %w", err)
	}
	return nil
}

func (ds *DataSet) AddSql(sql string) *DataSet {
	ds.Sql = append(ds.Sql, sql)

	return ds
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
				dataType:   DataType(Text),
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
