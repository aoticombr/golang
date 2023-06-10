package component

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Variant struct {
	Value interface{}
}
func (v Variant) SetValue(value interface{}){
	v.Value = value
}
func (v Variant) AsValue() interface{}{
  return	v.Value 
}

func (v Variant) AsString() string {
	switch val := v.Value.(type) {
	case time.Time:
		layout := "02/01/2006 15:04:05"
		formattedTime := val.Format(layout)
		if strings.Contains(formattedTime, "00:00:00") {
			formattedTime = formattedTime[:10]
		}
		return formattedTime
	case string:
		return val
	default:
		t := reflect.TypeOf(v.Value)
		panic(fmt.Sprintf("Unable to convert data type to string, Tipo: %v\n", t))
	}
}

func (v Variant) AsInt() int {
	switch val := v.Value.(type) {
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(val).Int())
	default:
		panic("Unable to convert data type to int")
	}
}

func (v Variant) AsInt64() int64 {
	switch val := v.Value.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(val).Int()
	case string:
		valueInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0
		}
		return valueInt
	default:
		return 0
	}
}

func (v Variant) AsFloat() float32 {
	switch val := v.Value.(type) {
	case float32:
		return val
	case float64:
		return float32(val)
	case string:
		floatValue, err := strconv.ParseFloat(val, 32)
		if err != nil {
			panic("Error converting string value to float32")
		}
		return float32(floatValue)
	default:
		t := reflect.TypeOf(val)
		panic(fmt.Sprintf("Unable to convert data type to float32, Tipo: %v\n", t))
	}
}

func (v Variant) AsFloat64() float64 {
	switch val := v.Value.(type) {
	case float32:
		return float64(val)
	case float64:
		return val
	case string:
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			panic("Error converting string value to float64")
		}
		return floatValue
	default:
		return 0
	}
}

func (v Variant) AsBool() bool {
	if v.Value != nil {
		switch val := v.Value.(type) {
		case int, int8, int16, int32, int64:
			return reflect.ValueOf(val).Int() == 1
		case string:
			v := strings.ToUpper(strings.TrimSpace(val))
			return v == "1" || v == "S" || v == "Y"
		default:
			panic("Unable to convert data type to bool")
		}
	}
	return false
}

func (v Variant) AsDateTime() time.Time {
	if v.Value != nil {
		switch val := v.Value.(type) {
		case time.Time:
			return val
		default:
			panic("Unable to convert data type to time.Time")
		}
	}
	return v.Value.(time.Time)
}