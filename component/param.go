package component

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Input int

const (
	IN    Input = 0
	OUT   Input = 1
	INOUT Input = 2
)

type Value interface{ *Variant | Variant }

type Param struct {
	Value Variant
	Input Input
	Tipo  reflect.Type
}

func (p Param) SetValue(value any) {
	p.Value.Value = value
}

func (p Param) asValue() Variant {
	if IsPointer(p.Value.Value) {
		fmt.Println("IsPointer")
		a := reflect.ValueOf(p.Value.Value).Elem().Interface()
		return Variant{Value: a}
	} else {
		fmt.Println("not IsPointer")
		return p.Value
	}

}

func (p Param) AsValue() interface{} {
	return p.asValue().AsValue()
}

func (p Param) AsString() string {
	return p.asValue().AsString()
}

func (p Param) AsInt() int {
	return p.asValue().AsInt()
}

func (p Param) AsInt64() int64 {
	return p.asValue().AsInt64()
}

func (p Param) AsFloat() float32 {
	return p.asValue().AsFloat()
}

func (p Param) AsFloat64() float64 {
	return p.asValue().AsFloat64()
}

func (p Param) AsBool() bool {
	return p.asValue().AsBool()
}

func (p Param) AsDateTime() time.Time {
	return p.asValue().AsDateTime()
}

type Params map[string]Param

func ConvertToInsertStatement(params Params) (string, string) {
	var columns []string
	var values []string
	for key, _ := range params {
		columns = append(columns, key)
		values = append(values, ":"+key)
	}
	columnsStr := strings.Join(columns, ", ")
	valuesStr := strings.Join(values, ", ")
	return columnsStr, valuesStr
}

//type Params map[string]Param
