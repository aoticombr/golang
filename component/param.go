package component

import (
	"fmt"
	"reflect"
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

func (p Param) asValue() Variant {
	//tp := reflect.TypeOf(p.Value.Value)
	//	fmt.Println("Param.TypeOf", tp)
	if IsPointer(p.Value.Value) {
		fmt.Println("IsPointer")
		a := reflect.ValueOf(p.Value.Value).Elem().Interface()
		tp := reflect.TypeOf(a)
		fmt.Println("Param.TypeOf", tp)
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

//type Params map[string]Param
