package github.com/aoticombr/godataset

import (
	"strings"
	"time"
)

type Field struct {
	name       string
	caption    string
	dataType   DataType
	Value      any
	dataMask   string
	valueTrue  string
	valueFalse string
	visible    bool
	order      int
	index      int
}

func (field Field) AsString() string {
	return field.Value.(string)
}

func (field Field) AsInt() int {
	if field.Value != nil {
		switch field.Value.(type) {
		case int:
			return field.Value.(int)
		case int8:
			return int(field.Value.(int8))
		case int16:
			return int(field.Value.(int16))
		case int32:
			return int(field.Value.(int32))
		case int64:
			return int(field.Value.(int64))
		default:
			panic("unable to convert data type to int")
		}
	} else {
		return 0
	}
}

func (field Field) AsInt64() int64 {
	if field.Value != nil {
		switch field.Value.(type) {
		case int:
			return int64(field.Value.(int))
		case int8:
			return int64(field.Value.(int8))
		case int16:
			return int64(field.Value.(int16))
		case int32:
			return int64(field.Value.(int32))
		case int64:
			return field.Value.(int64)
		default:
			panic("Unable to convert data type to int")
		}
	} else {
		return int64(0)
	}
}

func (field Field) AsFloat() float32 {
	if field.Value != nil {
		switch field.Value.(type) {
		case float32:
			return field.Value.(float32)
		case float64:
			return float32(field.Value.(float64))
		default:
			panic("Unable to convert data type to float32")
		}
	} else {
		return float32(0)
	}
}

func (field Field) AsFloat64() float64 {
	if field.Value != nil {
		switch field.Value.(type) {
		case float32:
			return float64(field.Value.(float32))
		case float64:
			return field.Value.(float64)
		default:
			panic("Unable to convert data type to float64")
		}
	} else {
		return float64(0)
	}
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
