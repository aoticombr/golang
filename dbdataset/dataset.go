package dbdataset

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"github.com/aoticombr/golang/dbconnbase"
	"github.com/aoticombr/golang/preprocesssql"
	"github.com/aoticombr/golang/stringlist"
	"github.com/aoticombr/golang/variant"
	goOra "github.com/sijms/go-ora/v2"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

type DS interface {
	NewDataSet(db *dbconnbase.Conn) *DataSet
	Open() error
	Close()
	Exec() (sql.Result, error)
	GetSql() (sql string)
	GetParams() []any
	Scan(list *sql.Rows)
	ParamByName(paramName string) Param
	SetInputParam(paramName string, paramValue any) *DataSet
	SetOutputParam(paramName string, paramType any) *DataSet
	FieldByName(fieldName string) Field
	Locate(key string, value any) bool
	First()
	Next()
	Eof() bool
	IsEmpty() bool
	IsNotEmpty() bool
	Count() int
	AddSql(sql string) *DataSet
	ToStruct(model any) error
	ToStructJson(model any) ([]byte, error)
}

type DataSet struct {
	Connection      *dbconnbase.Conn
	Transction      *dbconnbase.Transaction
	Ctx             context.Context
	Sql             stringlist.Strings
	Fields          *Fields
	Params          *Params
	Macros          *Macros
	Rows            []Row
	Index           int
	Recno           int
	IndexFieldNames string
	Silent          bool
	Name            string
	TagColumn       string
}

func newDataSetBase() *DataSet {
	ds := &DataSet{
		Index:  0,
		Recno:  0,
		Silent: true,
		Fields: NewFields(),
		Params: NewParams(),
		Macros: NewMacros(),
	}
	ds.Sql.Delimiter = " "
	ds.Fields.Owner = ds
	return ds
}

func NewDataSet(db *dbconnbase.Conn, opts ...OptionsDataset) *DataSet {
	ds := newDataSetBase()
	ds.Connection = db

	for _, opt := range opts {
		opt(ds)
	}

	return ds
}

func NewDataSetTx(tx *dbconnbase.Transaction, opts ...OptionsDataset) *DataSet {
	ds := newDataSetBase()
	ds.Transction = tx

	for _, opt := range opts {
		opt(ds)
	}

	return ds
}

func (ds *DataSet) AddContext(ctx context.Context) *DataSet {
	ds.Ctx = ctx
	return ds
}

func (ds *DataSet) Open() error {
	if ds.Connection != nil {
		if ds.Connection.Log {
			fmt.Printf("Open, name: %s\n", ds.Name)
		}
	}
	ds.Rows = nil
	ds.Index = 0
	ds.Recno = 0

	query := ds.getSql()

	var rows *sql.Rows
	var err error

	if ds.Transction != nil {
		if ds.Transction.Conn.Log {
			fmt.Println(query)
			ds.PrintParam()
		}

		if ds.Ctx != nil {
			rows, err = ds.Transction.Tx.QueryContext(ds.Ctx, query, ds.GetParams()...)
		} else {
			rows, err = ds.Transction.Tx.Query(query, ds.GetParams()...)
		}

		if err != nil {
			return err
		}
	} else {
		if ds.Connection.Log {
			fmt.Println(query)
			ds.PrintParam()
		}

		if ds.Ctx != nil {
			rows, err = ds.Connection.DB.QueryContext(ds.Ctx, query, ds.GetParams()...)
		} else {
			rows, err = ds.Connection.DB.Query(query, ds.GetParams()...)
		}

		if err != nil {
			errPing := ds.Connection.DB.Ping()
			if errPing != nil {
				errConn := ds.Connection.Open()
				if errConn != nil {
					return err
				}

				if ds.Ctx != nil {
					rows, err = ds.Connection.DB.QueryContext(ds.Ctx, query, ds.GetParams()...)
				} else {
					rows, err = ds.Connection.DB.Query(query, ds.GetParams()...)
				}

				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("could not open dataset %v", err)
			}
		}
	}

	if rows == nil {
		return fmt.Errorf("rows empty")
	}

	defer rows.Close()

	ds.scan(rows)
	ds.First()

	return nil
}

func (ds *DataSet) OpenContext(context context.Context) error {
	ds.Rows = nil
	ds.Index = 0
	ds.Recno = 0
	query := ds.getSql()

	if ds.Connection.Log {
		fmt.Println(query)
		ds.PrintParam()
	}

	var rows *sql.Rows
	var err error

	if ds.Transction != nil {
		if ds.Transction.Conn.Log {
			fmt.Println(query)
			ds.PrintParam()
		}

		rows, err = ds.Transction.Tx.QueryContext(context, query, ds.GetParams()...)

		if err != nil {
			return err
		}
	} else {
		if ds.Connection.Log {
			fmt.Println(query)
			ds.PrintParam()
		}

		rows, err = ds.Connection.DB.QueryContext(context, query, ds.GetParams()...)

		if err != nil {
			errPing := ds.Connection.DB.PingContext(context)
			if errPing != nil {
				errConn := ds.Connection.Open()
				if errConn != nil {
					return err
				}

				rows, err = ds.Connection.DB.QueryContext(context, query, ds.GetParams()...)

				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("could not open dataset %v", err)
			}
		}
	}

	if rows == nil {
		return fmt.Errorf("rows empty")
	}

	defer rows.Close()

	ds.scan(rows)

	ds.First()

	return nil
}

func (ds *DataSet) Close() {
	if ds.Connection != nil {
		if ds.Connection.Log {
			fmt.Printf("Close, name: %s\n", ds.Name)
		}
	}
	ds.Index = 0
	ds.Recno = 0
	ds.Rows = nil
	ds.Fields = nil
	ds.Params = nil
	ds.Macros = nil
	ds.IndexFieldNames = ""

	ds.Fields = NewFields()
	ds.Params = NewParams()
	ds.Macros = NewMacros()
	ds.Fields.Owner = ds

}

func (ds *DataSet) ClearParams() {
	if ds.Connection != nil {
		if ds.Connection.Log {
			fmt.Printf("ClearParams, name: %s\n", ds.Name)
		}
	}

	ds.Params = nil
	ds.Params = NewParams()

}
func (ds *DataSet) Free() {
	if ds.Connection != nil {
		if ds.Connection.Log {
			fmt.Printf("Free, name: %s\n", ds.Name)
		}
	}
	ds.Index = 0
	ds.Recno = 0
	ds.Rows = nil
	ds.Fields = nil
	ds.Params = nil
	ds.Macros = nil
	ds.IndexFieldNames = ""
	ds.Connection = nil

}

func (ds *DataSet) Exec() (sql.Result, error) {
	if ds.Connection != nil {
		if ds.Connection.Log {
			fmt.Printf("Exec, name: %s\n", ds.Name)
		}
	}

	query := ds.GetSql()

	if ds.Transction != nil {
		if ds.Transction.Conn.Log {
			fmt.Println(query)
			ds.PrintParam()
		}

		if ds.Ctx != nil {
			return ds.Transction.Tx.ExecContext(ds.Ctx, query, ds.GetParams()...)
		} else {
			return ds.Transction.Tx.Exec(query, ds.GetParams()...)
		}
	} else {
		if ds.Connection.Log {
			fmt.Println(query)
			ds.PrintParam()
		}

		if ds.Ctx != nil {
			return ds.Connection.DB.ExecContext(ds.Ctx, query, ds.GetParams()...)
		} else {
			return ds.Connection.DB.Exec(query, ds.GetParams()...)
		}
	}
}

func (ds *DataSet) ExecContext(context context.Context) (sql.Result, error) {
	if ds.Connection != nil {
		if ds.Connection.Log {
			fmt.Printf("ExecContext, name: %s\n", ds.Name)
		}
	}
	var stmt *sql.Stmt
	var err error

	vsql := ds.GetSql()

	if ds.Transction != nil {
		if ds.Transction.Conn.Log {
			fmt.Println(vsql)
			ds.PrintParam()
		}

		stmt, err = ds.Transction.Tx.PrepareContext(context, vsql)

		if err != nil {
			return nil, err
		}

	} else {
		if ds.Connection.Log {
			fmt.Println(vsql)
			ds.PrintParam()
		}

		stmt, err = ds.Connection.DB.PrepareContext(context, vsql)

		if err != nil {
			return nil, err
		}
	}

	defer stmt.Close()

	return stmt.ExecContext(context, ds.GetParams()...)
}

func (ds *DataSet) ExecBatch(size int) error {

	var stmt *sql.Stmt
	var err error

	query := ds.GetSql()

	if ds.Transction != nil {
		if ds.Transction.Conn.Log {
			fmt.Println(query)
			ds.PrintParam()
		}

		if ds.Ctx != nil {
			stmt, err = ds.Transction.Tx.PrepareContext(ds.Ctx, query)
		} else {
			stmt, err = ds.Transction.Tx.Prepare(query)
		}

		defer func() {
			_ = stmt.Close()
		}()

		if err != nil {
			return err
		}

	} else {
		if ds.Connection.Log {
			fmt.Println(query)
			ds.PrintParam()
		}

		if ds.Ctx != nil {
			stmt, err = ds.Connection.DB.PrepareContext(ds.Ctx, query)
		} else {
			stmt, err = ds.Connection.DB.Prepare(query)
		}

		defer func() {
			_ = stmt.Close()
		}()

		if err != nil {
			return err
		}
	}

	if ds.Ctx != nil {
		for i := 0; i < ds.Params.BatchSize; i++ {
			_, err = stmt.ExecContext(ds.Ctx, ds.GetParamsBatch(i)...)

			if err != nil {
				return err
			}
		}
	} else {
		for i := 0; i < ds.Params.BatchSize; i++ {
			_, err = stmt.Exec(ds.GetParamsBatch(i)...)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ds *DataSet) Delete() (int64, error) {

	var result sql.Result
	var err error

	if ds.Ctx != nil {
		result, err = ds.ExecContext(ds.Ctx)
	} else {
		result, err = ds.Exec()
	}

	if err != nil {
		return 0, err
	}

	rowsAff, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (ds *DataSet) DeleteContext(context context.Context) (int64, error) {
	result, err := ds.ExecContext(context)

	if err != nil {
		return 0, err
	}

	rowsAff, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsAff, nil
}

func (ds *DataSet) GetSql() (sql string) {

	sql = ds.Sql.Text()

	for i := 0; i < len(ds.Macros.List); i++ {
		key := ds.Macros.List[i].Name
		variant := ds.Macros.List[i].Value

		value := reflect.ValueOf(variant.Value)
		if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
			sql = strings.ReplaceAll(sql, "&"+key, JoinSlice(variant.AsValue()))
		} else {
			sql = strings.ReplaceAll(sql, "&"+key, variant.AsString())
		}
	}
	sql = strings.Replace(sql, "\r", " \n", -1)
	sql = strings.Replace(sql, "\n", " \n ", -1)
	sql = ds.replaceAllParam(sql)

	return sql
}

func (ds *DataSet) getSql() (vsql string) {

	vsql = ds.GetSql()

	vsql = strings.Replace(vsql, "\r", "\n", -1)
	vsql = strings.Replace(vsql, "\n", "\n ", -1)
	vsql = ds.replaceAllParam(vsql)

	return vsql
}

func (ds *DataSet) GetParams() []any {
	var param []any

	for i := 0; i < len(ds.Params.List); i++ {
		key := ds.Params.List[i].Name
		value := ds.Params.List[i].Value.Value

		switch ds.Params.List[i].ParamType {
		case IN:
			param = append(param, sql.Named(key, value))
		case OUT:
			param = append(param, sql.Named(key, sql.Out{Dest: value}))
		case INOUT:
			param = append(param, sql.Named(key, sql.Out{Dest: value, In: true}))
		}
	}
	return param
}

func (ds *DataSet) GetParamsBatch(index int) []any {
	var param []any

	for i := 0; i < len(ds.Params.List); i++ {
		key := ds.Params.List[i].Name
		value := ds.Params.List[i].Values[index].Value

		switch ds.Params.List[i].ParamType {
		case IN:
			param = append(param, sql.Named(key, value))
		case OUT:
			param = append(param, sql.Named(key, sql.Out{Dest: value}))
		case INOUT:
			param = append(param, sql.Named(key, sql.Out{Dest: value, In: true}))
		}
	}
	return param
}

func (ds *DataSet) GetMacros() []any {
	var macro []any
	for i := 0; i < len(ds.Macros.List); i++ {
		key := ds.Macros.List[i].Name
		value := &ds.Macros.List[i].Value.Value

		macro = append(macro, sql.Named(key, value))
	}
	return macro
}

type Lob struct {
	io.Reader
	IsClob bool
}

func (ds *DataSet) scan(list *sql.Rows) {
	fieldTypes, _ := list.ColumnTypes()
	fields, _ := list.Columns()

	if len(ds.Fields.List) == 0 {
		for i := 0; i < len(fields); i++ {
			field := ds.Fields.Add(fields[i])
			field.DataType = fieldTypes[i]
			switch field.DataType.ScanType().Kind() {
			case reflect.String:
				field.IDataType = new(DataType) //inicializa por é um ponteiro
				*field.IDataType = Text
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field.IDataType = new(DataType) //inicializa por é um ponteiro
				*field.IDataType = Integer
			case reflect.Float32, reflect.Float64:
				field.IDataType = new(DataType) //inicializa por é um ponteiro
				*field.IDataType = Float
			case reflect.Struct:
				if field.DataType.ScanType() == reflect.TypeOf(time.Time{}) {
					field.IDataType = new(DataType) //inicializa por é um ponteiro
					*field.IDataType = DateTime
				} else if field.DataType.ScanType() == reflect.TypeOf(&time.Time{}) {
					field.IDataType = new(DataType) //inicializa por é um ponteiro
					*field.IDataType = DateTime
				}
			case reflect.Bool:
				field.IDataType = new(DataType) //inicializa por é um ponteiro
				*field.IDataType = Boolean
			}
			field.Order = i + 1
			field.Index = i
		}
	}

	for list.Next() {
		columns := make([]any, len(fields))

		for i := 0; i < len(columns); i++ {
			columns[i] = &columns[i]
		}

		err := list.Scan(columns...)

		if err != nil {
			print(err)
		}

		row := NewRow()

		for i := 0; i < len(columns); i++ {
			row.List[strings.ToUpper(fields[i])] = &variant.Variant{
				Value: columns[i],
			}
		}

		ds.Rows = append(ds.Rows, row)
	}
}

func (ds *DataSet) ParamByName(paramName string) *Param {
	return ds.Params.ParamByName(paramName)
}

func (ds *DataSet) MacroByName(macroName string) *Macro {
	return ds.Macros.MacroByName(macroName)
}

func (ds *DataSet) SetInputParam(paramName string, paramValue any) *DataSet {
	ds.Params.SetInputParam(paramName, paramValue)
	return ds
}

func (ds *DataSet) SetInputParamClob(paramName string, paramValue string) *DataSet {
	ds.Params.SetInputParamClob(paramName, paramValue)
	return ds
}

func (ds *DataSet) SetInputParamBlob(paramName string, paramValue []byte) *DataSet {
	ds.Params.SetInputParamBlob(paramName, paramValue)
	return ds
}

func (ds *DataSet) SetOutputParam(paramName string, paramValue any) *DataSet {
	ds.Params.SetOutputParam(paramName, paramValue)
	return ds
}

func (ds *DataSet) SetOutputParamSlice(params ...ParamOut) *DataSet {
	ds.Params.SetOutputParamSlice(params...)
	return ds
}

func (ds *DataSet) SetMacro(macroName string, macroValue any) *DataSet {
	ds.Macros.SetMacro(macroName, macroValue)
	return ds
}

func (ds *DataSet) CreateFields() error {

	stmt, err := sqlparser.Parse(ds.GetSql())

	if err != nil {
		return fmt.Errorf("error when parsing the query %s, error: %w", ds.GetSql(), err)
	}

	sel, ok := stmt.(*sqlparser.Select)
	if ok {
		for _, expr := range sel.SelectExprs {
			alias, ok := expr.(*sqlparser.AliasedExpr)
			if ok {
				if !alias.As.IsEmpty() {
					_ = ds.Fields.Add(alias.As.String())
				} else {
					column, ok := alias.Expr.(*sqlparser.ColName)
					if ok {
						_ = ds.Fields.Add(column.Name.String())
					}
				}
			}
		}
	}

	return nil
}

func (ds *DataSet) Prepare() error {
	//Params, MacrosUpd, MacrosRead, err := preprocesssql.PreprocessSQL(ds.GetSql(), true, true, true, true, true)
	Params, _, _, err := preprocesssql.PreprocessSQL(ds.GetSql(), true, true, true, true, true)
	if err != nil {
		return err
	}
	for _, p := range Params.Items {
		param := &Param{
			Name:  p.Name,
			Value: &variant.Variant{Value: ""},
		}
		ds.Params.List = append(ds.Params.List, param)
	}
	return nil

}

func (ds *DataSet) FieldByName(fieldName string) *Field {
	return ds.Fields.FieldByName(fieldName)
}

func (ds *DataSet) Locate(key string, value any) bool {
	ds.First()
	for !ds.Eof() {
		switch value.(type) {
		case string:
			if ds.FieldByName(key).AsValue() == value {
				return true
			}
		default:
			if ds.FieldByName(key).AsValue() == value {
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
	if ds.Count() > 0 {
		ds.Recno = 1
	}
}

func (ds *DataSet) Next() {
	if !ds.Eof() {
		ds.Index++
		ds.Recno++
	}
}

func (ds *DataSet) Previous() {
	if !ds.Bof() {
		ds.Index--
		ds.Recno--
	}
}

func (ds *DataSet) Last() {
	ds.Index = ds.Count() - 1
	if ds.Count() == 0 {
		ds.Index = 0
	}
	ds.Recno = ds.Count()
}

func (ds *DataSet) Bof() bool {
	return ds.Count() == 0 || ds.Recno == 1
}

func (ds *DataSet) Eof() bool {
	return ds.Count() == 0 || ds.Recno > ds.Count()
}

func (ds *DataSet) IsEmpty() bool {
	return ds.Count() == 0
}

func (ds *DataSet) IsNotEmpty() bool {
	return ds.Count() > 0
}

func (ds *DataSet) Count() int {
	return len(ds.Rows)
}

func (ds *DataSet) AddSql(sql string) *DataSet {
	ds.Sql.Add(sql)

	return ds
}

func (ds *DataSet) ToStruct(model any) error {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	if modelType.Kind() == reflect.Pointer {
		modelType = modelValue.Type().Elem()
		modelValue = modelValue.Elem()
	} else {
		return fmt.Errorf("the variable %s is not a pointer", modelType.Name())
	}

	switch modelType.Kind() {
	case reflect.Struct:
		return ds.toStructUniqResult(modelValue)
	case reflect.Slice, reflect.Array:
		return ds.toStructList(modelValue)
	default:
		return errors.New("the interface is not a slice, array or struct")
	}
}

func (ds *DataSet) ToStructJson(model any) ([]byte, error) {
	err := ds.ToStruct(model)
	if err != nil {
		return nil, err
	}
	return json.Marshal(model)

}

func (ds *DataSet) toStructUniqResult(modelValue reflect.Value) error {
	for i := 0; i < modelValue.NumField(); i++ {
		if modelValue.Field(i).Kind() != reflect.Slice && modelValue.Field(i).Kind() != reflect.Array {
			field := modelValue.Type().Field(i)
			fieldName := field.Name
			fieldTag := field.Tag.Get(ds.TagColumn)
			fieldColumn := field.Tag.Get("column")
			fieldDb := field.Tag.Get("db")
			if fieldTag != "" && ds.TagColumn != "" {
				itens := strings.Split(fieldTag, ",")
				if strings.Contains(itens[0], "#") {
					continue
				}
				fieldName = itens[0]
			} else if fieldColumn != "" {
				itens := strings.Split(fieldColumn, ",")
				if strings.Contains(itens[0], "#") {
					continue
				}
				fieldName = itens[0]
			} else if fieldDb != "" {
				itens := strings.Split(fieldDb, ",")
				if strings.Contains(itens[0], "#") {
					continue
				}
				fieldName = itens[0]
			}

			if field.Anonymous {
				continue
			}

			var fieldValue any

			if field.Type.Kind() == reflect.Pointer {
				if ds.FieldByName(fieldName).IsNotNull() {
					fieldValue = ds.GetValue(ds.FieldByName(fieldName), modelValue.Field(i).Interface())

					if fieldValue != nil {
						rf := modelValue.Field(i)
						rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
						val := reflect.ValueOf(fieldValue)
						if rf.Kind() == reflect.Ptr && val.Kind() != reflect.Ptr {
							ptr := reflect.New(val.Type())
							ptr.Elem().Set(val)
							val = ptr
						}
						rf.Set(val)
					}
				}
			} else {
				fieldValue = ds.GetValue(ds.FieldByName(fieldName), modelValue.Field(i).Interface())

				if fieldValue != nil {
					rf := modelValue.Field(i)
					rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
					rf.Set(reflect.ValueOf(fieldValue))
				}
			}
		}
	}

	return nil
}

func (ds *DataSet) GetValue(field *Field, fieldType any) any {
	var fieldValue any
	valueType := reflect.TypeOf(fieldType)
	switch fieldType.(type) {
	case int:
		val := field.AsInt()
		fieldValue = reflect.ValueOf(val).Convert(valueType).Interface()
	case int8:
		val := field.AsInt8()
		fieldValue = reflect.ValueOf(val).Convert(valueType).Interface()
	case int16:
		val := field.AsInt16()
		fieldValue = reflect.ValueOf(val).Convert(valueType).Interface()
	case int32:
		val := field.AsInt32()
		fieldValue = reflect.ValueOf(val).Convert(valueType).Interface()
	case int64:
		val := field.AsInt64()
		fieldValue = reflect.ValueOf(val).Convert(valueType).Interface()
	case float32, float64:
		val := field.AsFloat64()
		fieldValue = reflect.ValueOf(val).Convert(valueType).Interface()
	case string:
		val := field.AsString()
		fieldValue = reflect.ValueOf(val).Convert(valueType).Interface()
	case bool:
		val := field.AsBool()
		fieldValue = reflect.ValueOf(val).Convert(valueType).Interface()
	case *int:
		val := field.AsInt()
		fieldValue = reflect.ValueOf(&val).Convert(valueType).Interface()
	case *int8:
		val := field.AsInt8()
		fieldValue = reflect.ValueOf(&val).Convert(valueType).Interface()
	case *int16:
		val := field.AsInt16()
		fieldValue = reflect.ValueOf(&val).Convert(valueType).Interface()
	case *int32:
		val := field.AsInt32()
		fieldValue = reflect.ValueOf(&val).Convert(valueType).Interface()
	case *int64:
		val := field.AsInt64()
		fieldValue = reflect.ValueOf(&val).Convert(valueType).Interface()
	case *float32, *float64:
		val := field.AsFloat64()
		fieldValue = reflect.ValueOf(&val).Convert(valueType).Interface()
	case *string:
		val := field.AsString()
		fieldValue = reflect.ValueOf(&val).Convert(valueType).Interface()
	case *bool:
		val := field.AsBool()
		fieldValue = reflect.ValueOf(&val).Convert(valueType).Interface()
	default:
		fieldValue = field.AsValue()
	}
	return fieldValue
}

func (ds *DataSet) toStructList(modelValue reflect.Value) error {
	var modelType reflect.Type

	if modelValue.Type().Elem().Kind() == reflect.Pointer {
		modelType = modelValue.Type().Elem()
	} else {
		modelType = modelValue.Type().Elem()
	}
	ds.First()
	for !ds.Eof() {
		newModel := reflect.New(modelType)

		if modelValue.Type().Elem().Kind() == reflect.Pointer {
			err := ds.toStructUniqResult(reflect.ValueOf(newModel.Interface()).Elem())
			if err != nil {
				return err
			}
		} else {
			err := ds.toStructUniqResult(reflect.ValueOf(newModel.Interface()).Elem())
			if err != nil {
				return err
			}
		}

		err := ds.toStructUniqResult(reflect.ValueOf(newModel.Interface()).Elem())
		if err != nil {
			return err
		}

		modelValue.Set(reflect.Append(modelValue, newModel.Elem()))

		ds.Next()
	}

	return nil
}

func JoinSlice(list any) string {
	valueOf := reflect.ValueOf(list)
	value := make([]string, valueOf.Len())
	for i := 0; i < valueOf.Len(); i++ {
		if valueOf.Index(i).Type().Kind() == reflect.String {
			value[i] = "'" + fmt.Sprintf("%v", valueOf.Index(i).Interface()) + "'"
		} else {
			value[i] = fmt.Sprintf("%v", valueOf.Index(i).Interface())
		}
	}
	return strings.Join(value, ", ")
}

func (ds *DataSet) PrintParam() {
	ds.Params.PrintParam()
}

func generateString() string {
	var result string
	for i := 0; i < 400; i++ {
		result += "abcdefghij"
	}
	return result
}

func (ds *DataSet) ParseSql() (sqlparser.Statement, error) {
	return sqlparser.Parse(ds.GetSql())
}

func limitStr(value string, limit int) string {
	if len(value) > limit {
		return value[:limit]
	}
	return value
}

func (ds *DataSet) SqlParam() string {
	vsql := ds.Sql.Text()

	replacements := make(map[string]string, len(ds.Params.List))

	for i := 0; i < len(ds.Params.List); i++ {
		key := ds.Params.List[i].Name
		value := ds.Params.List[i].Value

		var replacement string
		switch val := value.Value.(type) {
		case nil:
			replacement = "null"
		case time.Time:
			data := limitStr(value.AsString(), 19)
			replacement = "to_date('" + data + "','rrrr-mm-dd hh24:mi:ss')"
		case int:
			replacement = strconv.Itoa(val)
		case int8:
			replacement = strconv.FormatInt(int64(val), 10)
		case int16:
			replacement = strconv.FormatInt(int64(val), 10)
		case int32:
			replacement = strconv.FormatInt(int64(val), 10)
		case int64:
			replacement = strconv.FormatInt(val, 10)
		case uint:
			replacement = strconv.FormatUint(uint64(val), 10)
		case uint8:
			replacement = strconv.FormatUint(uint64(val), 10)
		case uint16:
			replacement = strconv.FormatUint(uint64(val), 10)
		case uint32:
			replacement = strconv.FormatUint(uint64(val), 10)
		case uint64:
			replacement = strconv.FormatUint(val, 10)
		case float32:
			replacement = strconv.FormatFloat(float64(val), 'f', -1, 32)
		case float64:
			replacement = strconv.FormatFloat(val, 'f', -1, 64)
		case string:
			replacement = "'" + val + "'"
		case []byte:
			replacement = "'" + string(val) + "'"
		case goOra.Clob:
			replacement = "q'[" + value.AsString() + "]'"
		}

		if ds.Params.List[i].ParamType != OUT {
			replacements[":"+key] = replacement
		}
	}

	// Faz as substituicoes em uma unica passada
	result := vsql
	for oldStr, newStr := range replacements {
		result = strings.ReplaceAll(result, oldStr, newStr)
	}

	return result
}

func StrNotEmpty(s string) bool {
	if len(s) == 0 {
		return false
	}

	r := []rune(s)
	l := len(r)

	for l > 0 {
		l--
		if !unicode.IsSpace(r[l]) {
			return true
		}
	}

	return false
}

func replaceParamPG(sql, param string, paramNumber int) (string, int) {
	pSize := len(param)

	i := strings.Index(sql, param)

	var ok bool
	switch {

	case i == -1:
		return sql, paramNumber

	case i == len(sql)-pSize:
		ok = true

	default:

		switch string(sql[i+pSize]) {
		case " ", ",", "(", ")", "=", "|", "[", "]":
			ok = true
		}
	}

	start := sql[:i]
	end, paramNumber := replaceParamPG(sql[i+pSize:], param, paramNumber)

	if ok {
		sql = fmt.Sprintf("%s$%d%s", start, paramNumber, end)
	} else {
		sql = start + param + end
	}

	return sql, paramNumber
}

func (ds *DataSet) replaceAllParam(sql string) (newSql string) {
	newSql = sql
	var dialect dbconnbase.DialectType

	if ds.Transction != nil {
		dialect = ds.Transction.Conn.Dialect
	} else {
		dialect = ds.Connection.Dialect
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error replaceParam", err)
			return
		}
	}()

	switch dialect {
	case dbconnbase.POSTGRESQL:
		for i := 0; i < len(ds.Params.List); i++ {
			param := ":" + ds.Params.List[i].Name
			newSql, _ = replaceParamPG(newSql, param, i+1)
		}
	}

	return
}
