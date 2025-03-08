package gorm

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

func GetTable(table interface{}) string {
	t := reflect.TypeOf(table).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("table")
		if fieldName != "" {
			return fieldName
		}
	}
	return ""
}
func GetPrimaryKey(table interface{}) []string {
	t := reflect.TypeOf(table)
	var primaryKeys []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("primarykey")
		if fieldName != "" {
			primaryKeys = append(primaryKeys, field.Name)
		}
	}
	return primaryKeys
}

type Column struct {
	Name        string
	PrimaryKey  bool
	UniqueKey   bool
	Required    bool
	Insert      bool
	Update      bool
	Delete      bool
	Where       bool
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
	return &Column{
		Name:        name,
		PrimaryKey:  false,
		UniqueKey:   false,
		Insert:      false,
		Update:      false,
		Delete:      false,
		Where:       false,
		Upper:       false,
		Lower:       false,
		ActionType:  false,
		Omitempty:   false,
		Nullempty:   false,
		Md5:         false,
		AutoGuid:    false,
		Required:    false,
		ReturnValue: false,
	}
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

func NewTable(table interface{}) *Table {
	tb := &Table{
		table: table,
		Options: Options{
			Delete: D_Remove,
			Db:     DB_Postgres,
		},
	}
	t := reflect.TypeOf(tb.table)
	v := reflect.ValueOf(tb.table)
	// Check if data is a pointer, if yes, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		Value := v.Field(i).Interface()
		table := field.Tag.Get("table")

		if table != "" {
			tb.TableName = table
		}

		column := field.Tag.Get("column")
		itens := strings.Split(column, ",")
		if len(itens) > 0 {

			if column != "" {
				col := NewColumn(itens[0])

				// Percorre e processa os itens
				for _, item := range itens {
					if item == "#actiontype" {
						col.ActionType = true
						switch Value.(string) {
						case "new":
							tb.ActionType = A_Insert
						case "old":
							tb.ActionType = A_Update
						case "del":
							tb.ActionType = A_Delete
						default:
							tb.ActionType = A_Insert
						}
						continue
					}
					if item == "#autoguid" {
						col.AutoGuid = true
					}
					if item == "#insert" {
						col.Insert = true
					}
					if item == "#update" {
						col.Update = true
					}
					if item == "#delete" {
						col.Delete = true
					}
					if item == "#where" {
						col.Delete = true
					}
					if item == "#primarykey" {
						col.PrimaryKey = true
					}
					if item == "#uniquekey" {
						col.UniqueKey = true
					}
					if item == "#required" {
						col.Required = true
					}
					if item == "#returnvalue" {
						col.ReturnValue = true
					}
					if item == "#omitempty" {
						col.Omitempty = true
					}
					if item == "#nullempty" {
						col.Nullempty = true
					}
					if item == "#md5" {
						col.Md5 = true
					}
					if item == "#upper" {
						col.Upper = true
					}
					if item == "#lower" {
						col.Lower = true
					}

				}
				tb.Columns = append(tb.Columns, col)

			}
		}
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
func (tb *Table) ValidateInc() error {

	t := reflect.TypeOf(tb.table)
	v := reflect.ValueOf(tb.table)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for _, Col := range tb.Columns {
		if Col.Insert {

			for b := 0; b < t.NumField(); b++ {
				field := t.Field(b)
				value := v.Field(b)

				column := field.Tag.Get("column")
				if idx := strings.Index(column, ","); idx != -1 {
					column = column[:idx]
				}
				if column != "" {
					if Col.Name == column {
						if Col.Required {
							switch value.Kind() {
							case reflect.Ptr:
								if value.IsNil() {
									return errors.New("column " + Col.Name + " is required")
								}
							case reflect.String:
								if value.String() == "" {
									return errors.New("column " + Col.Name + " is required")
								}
							}
						}

					}
				}

			}
		}
	}
	return nil

}

func (tb *Table) SqlInsert() (string, error) {
	err := tb.Validate()
	if err != nil {
		return "", err
	}
	err = tb.ValidateInc()
	if err != nil {
		return "", err
	}
	var columns, values, returncolumn, returninto string
	t := reflect.TypeOf(tb.table)
	v := reflect.ValueOf(tb.table)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for _, Col := range tb.Columns {

		if Col.Insert || Col.ReturnValue || Col.AutoGuid {

			for b := 0; b < t.NumField(); b++ {
				field := t.Field(b)
				value := v.Field(b)

				column := field.Tag.Get("column")
				if idx := strings.Index(column, ","); idx != -1 {
					column = column[:idx]
				}
				if column != "" && Col.ReturnValue {
					if Col.Name == column {
						if returncolumn != "" {
							returncolumn += ", "
							returninto += ", "
						}
						returncolumn += Col.Name
						returninto += ":new_" + Col.Name
					}

				}
				if column != "" && (Col.Insert || Col.AutoGuid) {
					if Col.Name == column {

						if Col.Omitempty {
							switch value.Kind() {
							case reflect.Ptr:
								if value.IsNil() {
									continue
								}
							case reflect.String:
								if value.String() == "" {
									continue
								}
							}
						}

						if columns != "" {
							columns += ", "
							values += ", "
						}

						columns += Col.Name
						if Col.Md5 {
							values += "md5("
						}
						if Col.Upper {
							values += "upper("
						}
						if Col.Lower {
							values += "lower("
						}
						if Col.Nullempty {
							switch value.Kind() {
							case reflect.Ptr:
								if value.IsNil() {
									values += "null"
								} else {
									values += ":" + Col.Name
								}
							}
						} else if Col.AutoGuid {
							switch value.Kind() {
							case reflect.Ptr:
								if value.IsNil() {
									values += "uuid_generate_v4()::uuid"
								} else {
									values += ":" + Col.Name
								}
							case reflect.String:
								if value.String() == "" {
									values += "uuid_generate_v4()::uuid"
								} else {
									values += ":" + Col.Name
								}
							default:
								values += ":" + Col.Name
							}

						} else {
							values += ":" + Col.Name
						}
						if Col.Upper {
							values += ")"
						}
						if Col.Lower {
							values += ")"
						}
						if Col.Md5 {
							values += ")"
						}
					}
				}

			}
		}
	}
	if columns == "" {
		return "", errors.New("columns not found")
	}
	if values == "" {
		return "", errors.New("values not found")
	}
	returnvalue := ""
	if tb.Columns.CountReturn() > 0 {
		switch tb.Options.Db {
		case DB_Oracle:
			returnvalue = " returning " + returncolumn + " into " + returninto
		case DB_Postgres:
			returnvalue = " returning " + returncolumn
		}
	}
	return "insert into " + tb.TableName + " (" + columns + ") values (" + values + ")" + returnvalue, nil
}

func (tb *Table) SqlUpdate() (string, error) {
	err := tb.Validate()
	if err != nil {
		return "", err
	}
	var (
		columns, where string
	)
	t := reflect.TypeOf(tb.table)
	v := reflect.ValueOf(tb.table)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for _, Col := range tb.Columns {
		if Col.Update {
			for b := 0; b < t.NumField(); b++ {
				field := t.Field(b)
				value := v.Field(b)

				column := field.Tag.Get("column")
				if idx := strings.Index(column, ","); idx != -1 {
					column = column[:idx]
				}
				if Col.Omitempty {
					switch value.Kind() {
					case reflect.Ptr:
						if value.IsNil() {
							continue
						}
					case reflect.String:
						if value.String() == "" {
							continue
						}
					}
				}
				if column != "" {
					if Col.Name == column {
						if columns != "" {
							columns += ", "
						}
						columns += Col.Name + "="

						if Col.Md5 {
							columns += "md5("
						}
						if Col.Upper {
							columns += "upper("
						}
						if Col.Lower {
							columns += "lower("
						}
						if Col.Nullempty {
							switch value.Kind() {
							case reflect.Ptr:
								if value.IsNil() {
									columns += "null"
								} else {
									columns += ":" + Col.Name
								}
							}
						} else {
							columns += ":" + Col.Name
						}

						if Col.Upper {
							columns += ")"
						}
						if Col.Lower {
							columns += ")"
						}
						if Col.Md5 {
							columns += ")"
						}
					}
				}
			}
		}
	}
	if tb.Options.Delete == D_Disable {
		where = "deleted_at is null"
	}
	for _, Col := range tb.Columns {
		if Col.PrimaryKey || Col.Where {
			if where != "" {
				where += " AND "
			}
			where += Col.Name + "=:" + Col.Name
		}
	}
	if columns == "" {
		return "", errors.New("columns not found")
	}
	return "UPDATE " + tb.TableName + " SET " + columns + " WHERE " + where, nil
}

func (tb *Table) SqlDelete() (string, error) {
	err := tb.Validate()
	if err != nil {
		return "", err
	}

	var where string
	if tb.Options.Delete == D_Disable {
		where = "deleted_at is null"
	}
	for _, Col := range tb.Columns {
		if Col.PrimaryKey || Col.Delete || Col.Where {

			if where != "" {
				where += " AND "
			}
			where += Col.Name + "=:" + Col.Name
		}

	}

	switch tb.Options.Delete {
	case D_Disable:
		return "UPDATE " + tb.TableName + " SET deleted_at=current_timestamp  WHERE " + where, nil
	case D_Remove:
		return "DELETE FROM " + tb.TableName + " WHERE " + where, nil
	}
	return "DELETE FROM " + tb.TableName + " WHERE " + where, nil
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
