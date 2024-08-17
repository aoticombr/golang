package gorm

import (
	"errors"
	"reflect"
	"strings"
)

type Delete int

const (
	D_Remove Delete = iota
	D_Disable
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
	Name      string
	Key       bool
	Insert    bool
	Update    bool
	Md5       bool
	Omitempty bool
}

func NewColumn(name string) *Column {
	return &Column{
		Name:   name,
		Key:    false,
		Insert: false,
		Update: false,
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
		if v.Key {
			count++
		}
	}
	return count
}

type Options struct {
	OmitColumnEmpty bool
	Delete          Delete
}

type Table struct {
	table     interface{}
	TableName string
	CRUD      string
	Columns   Columns
	Options   Options
}

func NewTable(table interface{}) *Table {
	tb := &Table{
		table: table,
		Options: Options{
			OmitColumnEmpty: false,
			Delete:          D_Remove,
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
		status := field.Tag.Get("json")
		itens := strings.Split(status, ",")
		// Percorre e processa os itens
		for _, item := range itens {
			if item == "crud" {
				//capturar valor string do campo
				tb.CRUD = Value.(string)
			}
		}

		if table != "" {
			tb.TableName = table
		}
		column := field.Tag.Get("column")
		itens = strings.Split(column, ",")
		if len(itens) > 0 {

			if column != "" {
				col := NewColumn(itens[0])

				// Percorre e processa os itens
				for _, item := range itens {
					if item == "insert" {
						col.Insert = true
					}
					if item == "update" {
						col.Update = true
					}
					if item == "primarykey" {
						col.Key = true
					}
					if item == "omitempty" {
						col.Omitempty = true
					}
					if item == "md5" {
						col.Md5 = true
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

func (tb *Table) SqlInsert() (string, error) {
	err := tb.Validate()
	if err != nil {
		return "", err
	}
	var columns, values string
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
						if tb.Options.OmitColumnEmpty {
							if value.Kind() == reflect.Ptr {
								// If it's a pointer and it's nil, skip it
								if value.IsNil() {
									continue
								}
							}
						}

						if columns != "" {
							columns += ", "
							values += ", "
						}
						if Col.Md5 {
							columns += "md5("
						}
						columns += Col.Name
						values += ":" + Col.Name
						if Col.Md5 {
							columns += ")"
						}
					}
				}

			}
		}
	}
	return "INSERT INTO " + tb.TableName + " (" + columns + ") VALUES (" + values + ")", nil
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
				if tb.Options.OmitColumnEmpty {
					if value.Kind() == reflect.Ptr {
						// If it's a pointer and it's nil, skip it
						if value.IsNil() {
							continue
						}
					}
				}
				if column != "" {
					if Col.Name == column {
						if columns != "" {
							columns += ", "
						}
						columns += Col.Name + "=:" + Col.Name
					}
				}
			}
		}
	}
	for _, Col := range tb.Columns {
		if Col.Key {
			if where != "" {
				where += " AND "
			}
			where += Col.Name + "=:" + Col.Name
		}
	}
	return "UPDATE " + tb.TableName + " SET " + columns + " WHERE " + where, nil
}

func (tb *Table) SqlDelete() (string, error) {
	err := tb.Validate()
	if err != nil {
		return "", err
	}

	var where string
	for _, Col := range tb.Columns {
		if Col.Key {
			if where != "" {
				where += " AND "
			}
			where += Col.Name + "=:" + Col.Name
		}

	}

	switch tb.Options.Delete {
	case D_Disable:
		return "UPDATE " + tb.TableName + " SET deleted_at=true WHERE " + where, nil
	case D_Remove:
		return "DELETE FROM " + tb.TableName + " WHERE " + where, nil
	}
	return "DELETE FROM " + tb.TableName + " WHERE " + where, nil
}

func (tb *Table) SqlStatus() (string, error) {
	switch tb.CRUD {
	case "new":
		return tb.SqlInsert()
	case "old":
		return tb.SqlUpdate()
	case "del":
		return tb.SqlDelete()
	}
	return "", errors.New("status not found")
}
