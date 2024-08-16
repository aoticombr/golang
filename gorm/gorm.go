package gorm

import (
	"reflect"
	"strings"
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

type Columns []string

func (c Columns) Exist(col string) bool {
	for _, v := range c {
		if v == col {
			return true
		}
	}
	return false
}

type Table struct {
	table       interface{}
	TableName   string
	PrimaryKeys Columns
	Columns     Columns
}

func NewTable(table interface{}) *Table {
	tb := &Table{
		table: table,
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
		table := field.Tag.Get("table")

		if table != "" {
			tb.TableName = table
		}
		column := field.Tag.Get("column")
		if idx := strings.Index(column, ","); idx != -1 {
			column = column[:idx]
		}
		if column != "" {
			tb.Columns = append(tb.Columns, column)
		}
		key := field.Tag.Get("primarykey")
		if key == "true" && column != "" {
			tb.PrimaryKeys = append(tb.PrimaryKeys, column)
		}

	}
	return tb
}

func (tb *Table) GetTableName() string {
	return tb.TableName
}

func (tb *Table) GetPrimaryKeys() Columns {
	return tb.PrimaryKeys
}

func (tb *Table) GetColumns() Columns {
	return tb.Columns
}

func (tb *Table) SqlInsert() (string, error) {
	var columns, values string
	for i := 0; i < len(tb.Columns); i++ {
		if i > 0 {
			columns += ", "
			values += ", "
		}
		columns += tb.Columns[i]
		values += ":" + tb.Columns[i]
	}
	return "INSERT INTO " + tb.TableName + " (" + columns + ") VALUES (" + values + ")", nil
}

func (tb *Table) SqlUpdate() (string, error) {
	var columns string
	for i := 0; i < len(tb.Columns); i++ {
		if !tb.PrimaryKeys.Exist(tb.Columns[i]) {
			if columns != "" {
				columns += ", "
			}
			columns += tb.Columns[i] + "=:" + tb.Columns[i]
		}
	}
	return "UPDATE " + tb.TableName + " SET " + columns, nil
}

func (tb *Table) SqlDelete() (string, error) {
	var columns string
	for i := 0; i < len(tb.PrimaryKeys); i++ {
		if i > 0 {
			columns += " AND "
		}
		columns += tb.PrimaryKeys[i] + "=:" + tb.PrimaryKeys[i]
	}
	return "DELETE FROM " + tb.TableName + " WHERE " + columns, nil
}

func (tb *Table) SqlInsertOmiteNil() (string, error) {
	var columns, values string
	t := reflect.TypeOf(tb.table)
	v := reflect.ValueOf(tb.table)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for a := 0; a < len(tb.Columns); a++ {
		for b := 0; b < t.NumField(); b++ {
			field := t.Field(b)
			value := v.Field(b)

			column := field.Tag.Get("column")
			if idx := strings.Index(column, ","); idx != -1 {
				column = column[:idx]
			}
			if column != "" {
				if tb.Columns[a] == column {
					if value.Kind() == reflect.Ptr {
						// If it's a pointer and it's nil, skip it
						if value.IsNil() {
							continue
						}
					}

					if columns != "" {
						columns += ", "
						values += ", "
					}
					columns += tb.Columns[a]
					values += ":" + tb.Columns[a]
				}
			}

		}
	}
	return "INSERT INTO " + tb.TableName + " (" + columns + ") VALUES (" + values + ")", nil
}
