package component

import (
	sql "database/sql"
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
	Value      any
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
		s := fmt.Sprintf("Unable to convert data type to string,Tipo: %v\n ", t)
		panic(s)
	}
}

func (field Field) AsInt() int {
	switch value := field.Value.(type) {
	case int:
		return int(value)
	case int8:
		return int(int8(value))
	case int16:
		return int(int16(value))
	case int32:
		return int(int32(value))
	case int64:
		return int(int64(value))
	default:
		panic("Unable to convert data type to int")
	}
}

func (field Field) AsInt64() int64 {
	switch value := field.Value.(type) {
	case int:
		return int64(field.Value.(int))
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case string:
		valueInt, err := convertStringToInt64(value)
		if err != nil {
			return 0
		}
		return valueInt
	}
	return 0
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
			panic("Erro ao converter valor string para float32")
		}
		return float32(floatValue)
	default:
		t := reflect.TypeOf(field.Value)
		s := fmt.Sprintf("Unable to convert data type to float32,Tipo: %v\n ", t)
		panic(s)
	}
}

func convertStringToInt64(value string) (int64, error) {
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return int64(number), nil
}

func (field Field) AsFloat64() float64 {

	switch value := field.Value.(type) {
	case float32:
		return float64(float32(value))
	case float64:
		return value
	case string:
		floatValue, err := strconv.ParseFloat(value, 32)
		if err != nil {
			panic("Erro ao converter valor string para float32")
		}
		return float64(floatValue)
	}
	return 0
}

func (field Field) AsBool() bool {
	if field.Value != nil {
		switch field.Value.(type) {
		case int:
			return field.Value.(int) == 1
		case int8:
			return field.Value.(int8) == 1
		case int16:
			return field.Value.(int16) == 1
		case int32:
			return field.Value.(int32) == 1
		case int64:
			return field.Value.(int64) == 1
		case string:
			value := strings.ToUpper(strings.Trim(field.Value.(string), " "))
			if value == "1" || value == "S" || value == "Y" {
				return true
			} else {
				return false
			}
		default:
			panic("Unable to convert data type to int")
		}
	} else {
		return false
	}

	return field.Value.(bool)
}

func (field Field) AsDateTime() time.Time {
	if field.Value != nil {
		switch field.Value.(type) {
		case time.Time:
			return field.Value.(time.Time)
		default:
			panic("Unable to convert data type to float64")
		}
	} else {
		data, _ := time.Parse(time.DateTime, time.DateTime)
		return data
	}
	return field.Value.(time.Time)
}

type Rows []map[string]Field
