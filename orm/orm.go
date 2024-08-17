package gorm

import (
	"errors"
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

type Options struct {
	OmitColumnEmpty bool
}

type Table struct {
	table       interface{}
	TableName   string
	CRUD        string
	PrimaryKeys Columns
	Columns     Columns
	Options     Options
}

func NewTable(table interface{}) *Table {
	tb := &Table{
		table: table,
		Options: Options{
			OmitColumnEmpty: false,
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

func (tb *Table) Validate() error {
	if tb.TableName == "" {
		return errors.New("table name not found")
	}
	if len(tb.PrimaryKeys) == 0 {
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
					columns += tb.Columns[a]
					values += ":" + tb.Columns[a]
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
	for a := 0; a < len(tb.Columns); a++ {
		if !tb.PrimaryKeys.Exist(tb.Columns[a]) {
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
					if tb.Columns[a] == column {
						if columns != "" {
							columns += ", "
						}
						columns += tb.Columns[a] + "=:" + tb.Columns[a]
					}
				}
			}
		}
	}
	for i := 0; i < len(tb.PrimaryKeys); i++ {
		if i > 0 {
			where += " AND "
		}
		where += tb.PrimaryKeys[i] + "=:" + tb.PrimaryKeys[i]
	}
	return "UPDATE " + tb.TableName + " SET " + columns + " WHERE " + where, nil
}

func (tb *Table) SqlDelete() (string, error) {
	err := tb.Validate()
	if err != nil {
		return "", err
	}

	var where string
	for i := 0; i < len(tb.PrimaryKeys); i++ {
		if i > 0 {
			where += " AND "
		}
		where += tb.PrimaryKeys[i] + "=:" + tb.PrimaryKeys[i]
	}
	if len(tb.PrimaryKeys) == 0 {
		return "", errors.New("primary key not found")
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
