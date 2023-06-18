package component

import (
	"reflect"
)

type Input int

const (
	IN    Input = 0
	OUT   Input = 1
	INOUT Input = 2
)

type Value interface{ *any | any }

type Param struct {
	Value Value
	Input Input
	Tipo  reflect.Type
}

func (p Param) GetData() any {

	value := p.Value.(*int64)
	switch p.Input {
	case INOUT:
		return *value
	default:
		return p.Value
	}

}

type Params map[string]Param

//type Params map[string]Param
