package dataset

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	cp "github.com/aoticombr/golang/component"
	conn "github.com/aoticombr/golang/connection"
)

type DataSet struct {
	Connection *conn.Conn
	Columns    []string
	Sql        cp.Strings
	rows       cp.Rows
	Params     cp.Params
	index      int
	Recno      int
	tx         *sql.Tx
}

func (ds *DataSet) Eof() bool {
	return ds.Count() == 0 || ds.Recno > ds.Count()
}
func (ds *DataSet) Count() int {
	return len(ds.rows)
}
func (ds *DataSet) GetParams() []any {
	var param []any
	for key, prm := range ds.Params {

		switch prm.Input {
		case cp.IN:
			param = append(param, sql.Named(key, prm.Value.Value))
		case cp.OUT:
			param = append(param, sql.Named(key, sql.Out{Dest: prm.Value.Value}))
		case cp.INOUT:
			param = append(param, sql.Named(key, sql.Out{Dest: prm.Value.Value, In: true}))
			//param = append(param, sql.Out{Dest: prm.Value, In: true})
		}
	}
	return param
}
func (ds *DataSet) Open() error {
	ds.rows = nil
	ds.index = 0
	ds.Recno = 0
	rows, err := ds.Connection.Db.Query(ds.Sql.Text(), ds.GetParams()...)

	if err != nil {
		return err
	}
	col, _ := rows.Columns()
	ds.Columns = col
	defer rows.Close()

	ds.scan(rows)

	ds.First()
	return nil
}

func (ds *DataSet) StartTransaction() error {
	if ds.tx != nil {
		t, err := ds.Connection.Db.Begin()
		if err != nil {
			return err
		}
		ds.tx = t
	}
	return nil
}
func (ds *DataSet) Commit() error {
	err := ds.tx.Commit()
	ds.tx = nil
	return err
}
func (ds *DataSet) Rollback() error {
	err := ds.tx.Rollback()
	ds.tx = nil
	return err
}
func (ds *DataSet) ExecTransact() (sql.Result, error) {
	return ds.tx.Exec(ds.Sql.Text(), ds.GetParams()...)
}
func (ds *DataSet) ExecDirect() (sql.Result, error) {
	stmt, err := ds.Connection.Db.Prepare(ds.Sql.Text())
	if err != nil {
		// handle the error
	}
	defer stmt.Close()
	return stmt.Exec(ds.GetParams()...)
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

func (ds *DataSet) SetInputParam(paramName string, paramValue any) *DataSet {

	ds.Params[paramName] = cp.Param{Value: cp.Variant{Value: paramValue}, Input: cp.IN}

	return ds
}

func (ds *DataSet) SetOutputParam(paramName string, tipo interface{}) *DataSet {

	switch tipo.(type) {
	case int, int8, int16, int32, int64:
		tipoValor := int64(0)
		ds.Params[paramName] = cp.Param{Value: cp.Variant{Value: &tipoValor}, Tipo: reflect.TypeOf(tipoValor), Input: cp.INOUT}
	case float32:
		tipoValor := float32(0)
		ds.Params[paramName] = cp.Param{Value: cp.Variant{Value: &tipoValor}, Tipo: reflect.TypeOf(tipoValor), Input: cp.INOUT}
	case float64:
		tipoValor := float64(0)
		ds.Params[paramName] = cp.Param{Value: cp.Variant{Value: &tipoValor}, Tipo: reflect.TypeOf(tipoValor), Input: cp.INOUT}
	case string:
		tipoValor := ""
		ds.Params[paramName] = cp.Param{Value: cp.Variant{Value: &tipoValor}, Tipo: reflect.TypeOf(tipoValor), Input: cp.INOUT}
	default:
		tipoValor := float64(0)
		ds.Params[paramName] = cp.Param{Value: cp.Variant{Value: &tipoValor}, Tipo: reflect.TypeOf(tipoValor), Input: cp.INOUT}
	}
	return ds
}
func (ds *DataSet) ParamByName(paramName string) cp.Param {
	return ds.Params[paramName]
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

func (ds *DataSet) RowInStructList(targetStruct interface{}) ([]interface{}, error) {

	targetType := reflect.TypeOf(targetStruct)
	if targetType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("targetStruct deve ser uma estrutura")
	}

	var results []interface{}
	fields := make(map[string]reflect.Value)

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		if field.Anonymous {
			continue // Ignorar campos anônimos
		}
		name := field.Name
		value := reflect.New(field.Type).Elem()
		fields[name] = value
	}

	for !ds.Eof() {
		result := reflect.New(targetType).Elem()

		for name, value := range fields {
			fieldValue := value.Interface()
			fieldType := value.Type()

			switch fieldValue.(type) {
			case int, int64, int32, int16, int8:
				val := ds.FieldByName(name).AsInt64()
				fieldValue = reflect.ValueOf(val).Convert(fieldType).Interface()
			case float32, float64:
				val := ds.FieldByName(name).AsFloat()
				fieldValue = reflect.ValueOf(val).Convert(fieldType).Interface()
			case string:
				val := ds.FieldByName(name).AsString()
				fieldValue = reflect.ValueOf(val).Convert(fieldType).Interface()
			case time.Time:
				val := ds.FieldByName(name).AsDateTime()
				fieldValue = reflect.ValueOf(val).Convert(fieldType).Interface()
			default:
				return nil, fmt.Errorf("tipo de campo não suportado: %v", fieldType)
			}

			result.FieldByName(name).Set(reflect.ValueOf(fieldValue))
		}

		results = append(results, result.Interface())
		ds.Next()
	}

	return results, nil
}

func (ds *DataSet) RowInStructObject(targetStruct interface{}) (interface{}, error) {
	results, err := ds.RowInStructList(targetStruct)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("nenhum objeto encontrado")
	}
	return results[0], nil
}

func GetDataSet(pconn *conn.Conn) *DataSet {
	ds := &DataSet{
		Connection: pconn,
		index:      0,
		Recno:      0,
		Params:     make(cp.Params),
	}
	return ds
}
