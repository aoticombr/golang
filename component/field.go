package component

import (
	"database/sql"
	"time"
)

type Field struct {
	Name       string
	Caption    string
	DataType   *sql.ColumnType
	Value      Variant
	DataMask   string
	ValueTrue  string
	ValueFalse string
	Visible    bool
	Order      int
	Index      int
}

func (field Field) AsValue() interface{} {
	return field.Value.AsValue()
}

func (field Field) AsString() string {
	return field.Value.AsString()
}

func (field Field) AsInt() int {
	return field.Value.AsInt()
}

func (field Field) AsInt64() int64 {
	return field.Value.AsInt64()
}

func (field Field) AsFloat() float32 {
	return field.Value.AsFloat()
}

func (field Field) AsFloat64() float64 {
	return field.Value.AsFloat64()
}

func (field Field) AsBool() bool {
	return field.Value.AsBool()
}

func (field Field) AsDateTime() time.Time {
	return field.Value.AsDateTime()
}

type Rows []map[string]Field
