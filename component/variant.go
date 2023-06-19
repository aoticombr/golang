package component

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Variant struct {
	Value any
}

func (v Variant) AsValue() any {
	return v.Value
}

func (v Variant) AsString() string {
	switch val := v.Value.(type) {
	case nil:
		return ""
	case time.Time:
		layout := "02/01/2006 15:04:05"
		formattedTime := val.Format(layout)
		if strings.Contains(formattedTime, "00:00:00") {
			formattedTime = formattedTime[:10]
		}
		return formattedTime
	case int, int8, int16, int32, int64:
		intValue := strconv.FormatInt(reflect.ValueOf(val).Int(), 10)
		return intValue
	case string:
		return val
	default:
		t := reflect.TypeOf(v.Value)
		log.Printf("unable to convert data type to string. Type: %v", t)
		return ""
	}
}

func (v Variant) AsInt() int {

	switch val := v.Value.(type) {
	case int:
		return v.Value.(int)
	case int8:
		return int(v.Value.(int8))
	case int16:
		return int(v.Value.(int16))
	case int32:
		return int(v.Value.(int32))
	case int64:
		return int(v.Value.(int64))
	case nil:
		return 0
	default:
		t := reflect.TypeOf(val)
		log.Printf("unable to convert data type to int, type: %v", t)
		return 0
	}

}

func (v Variant) AsInt64() int64 {

	switch val := v.Value.(type) {
	case int:
		return int64(v.Value.(int))
	case int8:
		return int64(v.Value.(int8))
	case int16:
		return int64(v.Value.(int16))
	case int32:
		return int64(v.Value.(int32))
	case int64:
		return v.Value.(int64)
	case nil:
		return int64(0)
	default:
		t := reflect.TypeOf(val)
		log.Printf("unable to convert data type to int64, type: %v", t)
		return int64(0)
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
			t := reflect.TypeOf(val)
			log.Printf("unable to convert data type to float32, type: %v", t)
			return float32(0)
		}
		return float32(floatValue)
	case nil:
		return float32(0)
	default:
		t := reflect.TypeOf(val)
		log.Printf("unable to convert data type to float32, type: %v", t)
		return float32(0)
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
			t := reflect.TypeOf(val)
			log.Printf("unable to convert data type to float64, type: %v", t)
			return float64(0)
		}
		return floatValue
	case nil:
		return float64(0)
	default:
		t := reflect.TypeOf(val)
		log.Printf("unable to convert data type to float64, type: %v", t)
		return float64(0)
	}

}

func (v Variant) AsBool() bool {

	switch val := v.Value.(type) {
	case int:
		return v.Value.(int) == 1
	case int8:
		return v.Value.(int8) == 1
	case int16:
		return v.Value.(int16) == 1
	case int32:
		return v.Value.(int32) == 1
	case int64:
		return v.Value.(int64) == 1
	case string:
		value := strings.ToUpper(strings.Trim(v.Value.(string), " "))
		if value == "1" || value == "S" || value == "Y" {
			return true
		} else {
			return false
		}
	case nil:
		return false
	default:
		t := reflect.TypeOf(val)
		log.Printf("unable to convert data type to bool, type: %v", t)
		return false
	}

}

func (v Variant) AsDateTime() time.Time {

	switch v.Value.(type) {
	case time.Time:
		return v.Value.(time.Time)
	case nil:
		data, _ := time.Parse(time.DateTime, time.DateTime)
		return data
	default:
		log.Printf("unable to convert data type to time.")
		data, _ := time.Parse(time.DateTime, time.DateTime)
		return data
	}

}

func IsPointer(value interface{}) bool {
	// Obtenha o tipo da variável
	t := reflect.TypeOf(value)

	// Verifique se o tipo é um ponteiro
	return t.Kind() == reflect.Ptr
}
