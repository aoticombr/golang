package orm

import (
	"errors"
	"reflect"
	"strings"
)

type DB int

const (
	DB_Oracle DB = iota
	DB_Postgres
)

type Delete int

const (
	D_Remove Delete = iota
	D_Disable
)

type ActionType int

const (
	A_Insert ActionType = iota
	A_Update
	A_Delete
)

func structType(table interface{}) reflect.Type {
	t := reflect.TypeOf(table)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func GetTable(table interface{}) string {
	t := structType(table)
	for i := 0; i < t.NumField(); i++ {
		if name := t.Field(i).Tag.Get("table"); name != "" {
			return name
		}
	}
	return ""
}

func GetPrimaryKey(table interface{}) []string {
	t := structType(table)
	var primaryKeys []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("primarykey") != "" {
			primaryKeys = append(primaryKeys, field.Name)
		}
	}
	return primaryKeys
}

type Column struct {
	Name        string
	fieldIndex  int
	PrimaryKey  bool
	UniqueKey   bool
	Required    bool
	Insert      bool
	Update      bool
	Delete      bool
	Where       bool
	TimeNow     bool
	Md5         bool
	Upper       bool
	Lower       bool
	AutoGuid    bool
	Omitempty   bool
	Nullempty   bool
	ActionType  bool
	ReturnValue bool
}

func NewColumn(name string) *Column {
	return &Column{Name: name, fieldIndex: -1}
}

type Columns []*Column

func (c Columns) Exist(col string) bool {
	for _, v := range c {
		if v.Name == col {
			return true
		}
	}
	return false
}

func (c Columns) CountKeys() int {
	var count int
	for _, v := range c {
		if v.PrimaryKey {
			count++
		}
	}
	return count
}

func (c Columns) CountReturn() int {
	var count int
	for _, v := range c {
		if v.ReturnValue {
			count++
		}
	}
	return count
}

type Options struct {
	Delete Delete
	Db     DB
}

type Table struct {
	table      interface{}
	TableName  string
	ActionType ActionType
	Columns    Columns
	Options    Options
}

// applyFlag aplica um item da tag em col. Aceita "#flag" e "flag" (retrocompat).
// Retorna false quando o item não corresponde a nenhuma flag conhecida.
func applyFlag(col *Column, item string) bool {
	switch strings.TrimPrefix(item, "#") {
	case "autoguid":
		col.AutoGuid = true
	case "insert":
		col.Insert = true
	case "update":
		col.Update = true
	case "delete":
		col.Delete = true
	case "where":
		col.Where = true
	case "timenow":
		col.TimeNow = true
	case "primarykey":
		col.PrimaryKey = true
	case "uniquekey":
		col.UniqueKey = true
	case "required":
		col.Required = true
	case "returnvalue":
		col.ReturnValue = true
	case "omitempty":
		col.Omitempty = true
	case "nullempty":
		col.Nullempty = true
	case "md5":
		col.Md5 = true
	case "upper":
		col.Upper = true
	case "lower":
		col.Lower = true
	case "actiontype":
		col.ActionType = true
	default:
		return false
	}
	return true
}

func NewTable(table interface{}) *Table {
	tb := &Table{
		table: table,
		Options: Options{
			Delete: D_Remove,
			Db:     DB_Postgres,
		},
	}
	t := reflect.TypeOf(table)
	v := reflect.ValueOf(table)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if name := field.Tag.Get("table"); name != "" {
			tb.TableName = name
		}

		column := field.Tag.Get("column")
		if column == "" || column == "-" {
			continue
		}
		itens := strings.Split(column, ",")
		colName := strings.TrimSpace(itens[0])
		if colName == "" {
			continue
		}
		col := NewColumn(colName)
		col.fieldIndex = i

		for _, item := range itens[1:] {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			applyFlag(col, item)

			if col.ActionType && strings.TrimPrefix(item, "#") == "actiontype" {
				s, _ := v.Field(i).Interface().(string)
				switch s {
				case "old":
					tb.ActionType = A_Update
				case "del":
					tb.ActionType = A_Delete
				default:
					tb.ActionType = A_Insert
				}
			}
		}
		tb.Columns = append(tb.Columns, col)
	}
	return tb
}

func (tb *Table) GetTableName() string {
	return tb.TableName
}

func (tb *Table) GetColumns() Columns {
	return tb.Columns
}

func (tb *Table) Validate() error {
	if tb.TableName == "" {
		return errors.New("table name not found")
	}
	if tb.Columns.CountKeys() == 0 {
		return errors.New("primary key not found")
	}
	if len(tb.Columns) == 0 {
		return errors.New("columns not found")
	}
	return nil
}

func (tb *Table) fieldValue(col *Column) (reflect.Value, bool) {
	v := reflect.ValueOf(tb.table)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if col.fieldIndex < 0 || col.fieldIndex >= v.NumField() {
		return reflect.Value{}, false
	}
	return v.Field(col.fieldIndex), true
}

func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
		return value.IsNil()
	case reflect.String:
		return value.String() == ""
	}
	return false
}

func (tb *Table) ValidateRequired() error {
	for _, col := range tb.Columns {
		if !col.Insert || !col.Required {
			continue
		}
		value, ok := tb.fieldValue(col)
		if !ok {
			continue
		}
		if isEmpty(value) {
			return errors.New("column " + col.Name + " is required")
		}
	}
	return nil
}

func wrapValue(col *Column, base string) string {
	if col.Md5 {
		base = "md5(" + base + ")"
	}
	if col.Upper {
		base = "upper(" + base + ")"
	}
	if col.Lower {
		base = "lower(" + base + ")"
	}
	return base
}

func (tb *Table) autoGuidExpr() string {
	if tb.Options.Db == DB_Oracle {
		return "sys_guid()"
	}
	return "uuid_generate_v4()::uuid"
}

func (tb *Table) SqlInsert() (string, error) {
	if err := tb.Validate(); err != nil {
		return "", err
	}
	if err := tb.ValidateRequired(); err != nil {
		return "", err
	}

	var columns, values, returncolumn, returninto []string
	for _, col := range tb.Columns {
		if col.ReturnValue {
			returncolumn = append(returncolumn, col.Name)
			returninto = append(returninto, ":new_"+col.Name)
		}

		if !(col.Insert || col.AutoGuid || col.TimeNow) {
			continue
		}
		value, ok := tb.fieldValue(col)
		if !ok {
			continue
		}
		if col.Omitempty && !col.TimeNow && isEmpty(value) {
			continue
		}

		var expr string
		switch {
		case col.TimeNow && col.Insert:
			expr = "current_timestamp"
		case col.Nullempty && isEmpty(value):
			expr = "null"
		case col.AutoGuid && isEmpty(value):
			expr = tb.autoGuidExpr()
		default:
			expr = ":" + col.Name
		}

		columns = append(columns, col.Name)
		values = append(values, wrapValue(col, expr))
	}

	if len(columns) == 0 {
		return "", errors.New("columns not found")
	}
	if len(values) == 0 {
		return "", errors.New("values not found")
	}

	sql := "insert into " + tb.TableName +
		" (" + strings.Join(columns, ", ") + ")" +
		" values (" + strings.Join(values, ", ") + ")"

	if len(returncolumn) > 0 {
		switch tb.Options.Db {
		case DB_Oracle:
			sql += " returning " + strings.Join(returncolumn, ", ") +
				" into " + strings.Join(returninto, ", ")
		case DB_Postgres:
			sql += " returning " + strings.Join(returncolumn, ", ")
		}
	}
	return sql, nil
}

func (tb *Table) SqlUpdate() (string, error) {
	if err := tb.Validate(); err != nil {
		return "", err
	}
	if err := tb.ValidateRequired(); err != nil {
		return "", err
	}

	var sets, where []string
	for _, col := range tb.Columns {
		if !col.Update {
			continue
		}
		value, ok := tb.fieldValue(col)
		if !ok {
			continue
		}
		if col.Omitempty && !col.TimeNow && isEmpty(value) {
			continue
		}

		var expr string
		switch {
		case col.TimeNow && col.Update:
			expr = "current_timestamp"
		case col.Nullempty && isEmpty(value):
			expr = "null"
		default:
			expr = ":" + col.Name
		}
		sets = append(sets, col.Name+" = "+wrapValue(col, expr))
	}

	if len(sets) == 0 {
		return "", errors.New("columns not found")
	}

	if tb.Options.Delete == D_Disable {
		where = append(where, "deleted_at is null")
	}
	for _, col := range tb.Columns {
		if col.PrimaryKey || col.Where {
			where = append(where, col.Name+"=:"+col.Name)
		}
	}
	return "UPDATE " + tb.TableName +
		" SET " + strings.Join(sets, ", ") +
		" WHERE " + strings.Join(where, " AND "), nil
}

func (tb *Table) SqlDelete() (string, error) {
	if err := tb.Validate(); err != nil {
		return "", err
	}

	var where []string
	if tb.Options.Delete == D_Disable {
		where = append(where, "deleted_at is null")
	}
	for _, col := range tb.Columns {
		if col.PrimaryKey || col.Delete || col.Where {
			where = append(where, col.Name+"=:"+col.Name)
		}
	}

	switch tb.Options.Delete {
	case D_Disable:
		return "UPDATE " + tb.TableName +
			" SET deleted_at = current_timestamp" +
			" WHERE " + strings.Join(where, " AND "), nil
	default:
		return "DELETE FROM " + tb.TableName +
			" WHERE " + strings.Join(where, " AND "), nil
	}
}

func (tb *Table) SqlStatus() (string, error) {
	switch tb.ActionType {
	case A_Insert:
		return tb.SqlInsert()
	case A_Update:
		return tb.SqlUpdate()
	case A_Delete:
		return tb.SqlDelete()
	}
	return "", errors.New("status not found")
}
