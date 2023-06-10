package component

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func (field Field) AsString() string {
	switch fieldValue := field.Value.(type) {
	case time.Time:
		layout := "02/01/2006 15:04:05"
		formattedTime := fieldValue.Format(layout)
		if strings.Contains(formattedTime, "00:00:00") {
			formattedTime = formattedTime[:10]
		}
		return formattedTime
	case string:
		return fieldValue
	default:
		t := reflect.TypeOf(field.Value)
		panic(fmt.Sprintf("Unable to convert data type to string, Tipo: %v\n", t))
	}
}

func (field Field) AsInt() int {
	switch value := field.Value.(type) {
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(value).Int())
	default:
		panic("Unable to convert data type to int")
	}
}

func (field Field) AsInt64() int64 {
	switch value := field.Value.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(value).Int()
	case string:
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0
		}
		return valueInt
	default:
		return 0
	}
}

func (field Field) AsFloat() float32 {
	switch value := field.Value.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case string:
		floatValue, err := strconv.ParseFloat(value, 32)
		if err != nil {
			panic("Error converting string value to float32")
		}
		return float32(floatValue)
	default:
		t := reflect.TypeOf(field.Value)
		panic(fmt.Sprintf("Unable to convert data type to float32, Tipo: %v\n", t))
	}
}

func (field Field) AsFloat64() float64 {
	switch value := field.Value.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case string:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Error converting string value to float64")
		}
		return floatValue
	default:
		return 0
	}
}

func (field Field) AsBool() bool {
	if field.Value != nil {
		switch value := field.Value.(type) {
		case int, int8, int16, int32, int64:
			return reflect.ValueOf(value).Int() == 1
		case string:
			v := strings.ToUpper(strings.TrimSpace(value))
			return v == "1" || v == "S" || v == "Y"
		default:
			panic("Unable to convert data type to bool")
		}
	}
	return false
}

func (field Field) AsDateTime() time.Time {
	if field.Value != nil {
		switch value := field.Value.(type) {
		case time.Time:
			return value
		default:
			panic("Unable to convert data type to time.Time")
		}
	}
	return field.Value.(time.Time)
}

type Rows []map[string]Field
