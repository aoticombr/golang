package variant

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	go_ora "github.com/sijms/go-ora/v2"
)

type Variant struct {
	Value  any
	Silent bool
}

func (v *Variant) SetSilent(value bool) *Variant {
	v.Silent = value
	return v
}

func (v *Variant) AsValue() any {
	return v.Value
}

func (v *Variant) AsString() string {
	value := ""
	switch val := v.Value.(type) {
	case nil:
		value = ""
	case time.Time:
		value = val.String()
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		value = fmt.Sprintf("%v", val)
	case float32, float64:
		value = fmt.Sprintf("%f", val)
	case string:
		value = val
	case []uint8:
		value = string([]byte(val))
	default:
		t := reflect.TypeOf(v.Value)
		msg := fmt.Sprintf("unable to convert data type to string. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}

		value = ""
	}
	value = strings.Replace(value, "\r", "\n", -1)
	return value
}
func (v *Variant) AsStringNil() *string {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsString()
		t, ok := tvalor.(string)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsInt() int {
	switch val := v.Value.(type) {
	case nil:
		return 0
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
	case uint:
		return int(v.Value.(uint))
	case uint8:
		return int(v.Value.(uint8))
	case uint16:
		return int(v.Value.(uint16))
	case uint32:
		return int(v.Value.(uint32))
	case uint64:
		return int(v.Value.(uint64))
	case string:
		intValue, err := strconv.Atoi(val)
		if err != nil {
			t := reflect.TypeOf(val)
			msg := fmt.Sprintf("unable to convert data type to int. Type: %v", t)
			if v.Silent {
				fmt.Println(msg)
			} else {
				panic(msg)
			}
			return 0
		}
		return intValue
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to int. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return 0
	}
}
func (v *Variant) AsIntNil() *int {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsInt()
		t, ok := tvalor.(int)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsInt8() int8 {
	switch val := v.Value.(type) {
	case nil:
		return int8(0)
	case int:
		return int8(val)
	case int8:
		return val
	case int16:
		return int8(val)
	case int32:
		return int8(val)
	case int64:
		return int8(val)
	case uint:
		return int8(val)
	case uint8:
		return int8(val)
	case uint16:
		return int8(val)
	case uint32:
		return int8(val)
	case uint64:
		return int8(val)
	case string:
		int8Value, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			t := reflect.TypeOf(val)
			msg := fmt.Sprintf("unable to convert data type to int8. Type: %v", t)
			if v.Silent {
				fmt.Println(msg)
			} else {
				panic(msg)
			}
			return int8(0)
		}
		return int8(int8Value)
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to int8. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return int8(0)
	}
}
func (v *Variant) AsInt8Nil() *int8 {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsInt8()
		t, ok := tvalor.(int8)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsInt16() int16 {
	switch val := v.Value.(type) {
	case nil:
		return int16(0)
	case int:
		return int16(val)
	case int8:
		return int16(val)
	case int16:
		return val
	case int32:
		return int16(val)
	case int64:
		return int16(val)
	case uint:
		return int16(val)
	case uint8:
		return int16(val)
	case uint16:
		return int16(val)
	case uint32:
		return int16(val)
	case uint64:
		return int16(val)
	case string:
		int16Value, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			t := reflect.TypeOf(val)
			msg := fmt.Sprintf("unable to convert data type to int16. Type: %v", t)
			if v.Silent {
				fmt.Println(msg)
			} else {
				panic(msg)
			}
			return int16(0)
		}
		return int16(int16Value)
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to int16. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return int16(0)
	}
}
func (v *Variant) AsInt16Nil() *int16 {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsInt16()
		t, ok := tvalor.(int16)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsInt32() int32 {
	switch val := v.Value.(type) {
	case nil:
		return int32(0)
	case int:
		return int32(val)
	case int8:
		return int32(val)
	case int16:
		return int32(val)
	case int32:
		return val
	case int64:
		return int32(val)
	case uint:
		return int32(val)
	case uint8:
		return int32(val)
	case uint16:
		return int32(val)
	case uint32:
		return int32(val)
	case uint64:
		return int32(val)
	case float32:
		return int32(val)
	case float64:
		return int32(val)
	case string:
		int32Value, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			t := reflect.TypeOf(val)
			msg := fmt.Sprintf("unable to convert data type to int32. Type: %v", t)
			if v.Silent {
				fmt.Println(msg)
			} else {
				panic(msg)
			}
			return int32(0)
		}
		return int32(int32Value)
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to int32. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return int32(0)
	}
}
func (v *Variant) AsInt32Nil() *int32 {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsInt32()
		t, ok := tvalor.(int32)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsInt64() int64 {
	switch val := v.Value.(type) {
	case nil:
		return int64(0)
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	case uint:
		return int64(val)
	case uint8:
		return int64(val)
	case uint16:
		return int64(val)
	case uint32:
		return int64(val)
	case uint64:
		return int64(val)
	case float32:
		return int64(val)
	case float64:
		return int64(val)
	case string:
		int64Value, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			t := reflect.TypeOf(val)
			msg := fmt.Sprintf("unable to convert data type to int64. Type: %v", t)
			if v.Silent {
				fmt.Println(msg)
			} else {
				panic(msg)
			}
			return int64(0)
		}
		return int64Value
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to int64. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return int64(0)
	}
}
func (v *Variant) AsInt64Nil() *int64 {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsInt64()
		t, ok := tvalor.(int64)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsFloat() float32 {
	switch val := v.Value.(type) {
	case nil:
		return float32(0)
	case int:
		return float32(val)
	case int8:
		return float32(val)
	case int16:
		return float32(val)
	case int32:
		return float32(val)
	case int64:
		return float32(val)
	case uint:
		return float32(val)
	case uint8:
		return float32(val)
	case uint16:
		return float32(val)
	case uint32:
		return float32(val)
	case uint64:
		return float32(val)
	case float32:
		return val
	case float64:
		return float32(val)
	case string:
		floatValue, err := strconv.ParseFloat(val, 32)
		if err != nil {
			t := reflect.TypeOf(val)
			msg := fmt.Sprintf("unable to convert data type to float32. Type: %v", t)
			if v.Silent {
				fmt.Println(msg)
			} else {
				panic(msg)
			}
			return float32(0)
		}
		return float32(floatValue)
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to float32. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return float32(0)
	}
}
func (v *Variant) AsFloatNil() *float32 {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsFloat()
		t, ok := tvalor.(float32)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsFloat64() float64 {
	switch val := v.Value.(type) {
	case nil:
		return float64(0)
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	case string:
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			t := reflect.TypeOf(val)
			msg := fmt.Sprintf("unable to convert data type to float64. Type: %v", t)
			if v.Silent {
				fmt.Println(msg)
			} else {
				panic(msg)
			}
			return float64(0)
		}
		return floatValue
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to float64. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return float64(0)
	}
}
func (v *Variant) AsFloat64Nil() *float64 {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsFloat64()
		t, ok := tvalor.(float64)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsBool() bool {
	switch val := v.Value.(type) {
	case nil:
		return false
	case bool:
		return v.Value.(bool)
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
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to bool. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return false
	}
}
func (v *Variant) AsBoolNil() *bool {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsBool()
		t, ok := tvalor.(bool)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) AsDateTime() time.Time {
	switch val := v.Value.(type) {
	case nil:
		data, _ := time.Parse(time.DateTime, time.DateTime)
		return data
	case time.Time:
		return v.Value.(time.Time)
	case string:
		dataLocal, err := dateparse.ParseAny(val)

		if err != nil {
			preferMonthFirstFalse := dateparse.PreferMonthFirst(false)
			dataLocal, err = dateparse.ParseAny(val, preferMonthFirstFalse)

			if err != nil {
				msg := fmt.Sprintf("unable to convert data type to time.")
				if v.Silent {
					fmt.Println(msg)
				} else {
					panic(msg)
				}
				data, _ := time.Parse(time.DateTime, time.DateTime)
				return data
			}
		}

		return dataLocal
	default:
		msg := fmt.Sprintf("unable to convert data type to time. ")
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		data, _ := time.Parse(time.DateTime, time.DateTime)
		return data
	}
}

func (v *Variant) AsDateTimeNil() *time.Time {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsDateTime()
		t, ok := tvalor.(time.Time)
		if ok {
			return &t
		}
		return nil
	}
}

func (v *Variant) AsByte() []byte {
	switch val := v.Value.(type) {
	case nil:
		return nil
	case []byte:
		return v.Value.([]byte)
	case string:
		return []byte(v.AsString())
	case *go_ora.Clob:
		return []byte(val.String)
	default:
		t := reflect.TypeOf(val)
		msg := fmt.Sprintf("unable to convert data type to byte. Type: %v", t)
		if v.Silent {
			fmt.Println(msg)
		} else {
			panic(msg)
		}
		return nil
	}
}

func (v *Variant) AsByteNil() *[]byte {
	valor := v.AsValue()
	var tvalor any

	if valor == nil {
		return nil
	} else {
		tvalor = v.AsByte()
		t, ok := tvalor.([]byte)
		if ok {
			return &t
		}
		return nil
	}
}
func (v *Variant) IsNull() bool {
	switch val := v.Value.(type) {
	case nil:
		return true
	case time.Time:
		return val.IsZero()
	case *time.Time:
		return val == nil
	case string:
		return val == ""
	default:
		return false
	}
}

func (v *Variant) IsNotNull() bool {
	return !v.IsNull()
}
