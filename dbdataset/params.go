package dbdataset

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	dbconnect "github.com/aoticombr/golang/dbconnbase"
	"github.com/aoticombr/golang/variant"
	goOra "github.com/sijms/go-ora/v2"
)

type Params struct {
	Owner     *DataSet
	BatchSize int
	List      []*Param
}

func NewParams() *Params {
	value := &Params{
		List: []*Param{},
	}
	return value
}

func (p *Params) FindParamByName(paramName string) *Param {
	var param *Param
	for i := 0; i < len(p.List); i++ {
		if strings.ToUpper(p.List[i].Name) == strings.ToUpper(paramName) {
			param = p.List[i]
		}
	}

	return param
}

func (p *Params) ParamByName(paramName string) *Param {
	var param *Param
	for i := 0; i < len(p.List); i++ {
		if strings.ToUpper(p.List[i].Name) == strings.ToUpper(paramName) {
			param = p.List[i]
		}
	}

	if param == nil {
		param = &Param{}
		fmt.Println("Parameter " + paramName + " doesn't exists")
	}

	return param
}

func (p *Params) Add(paramName string) *Param {
	param := p.FindParamByName(paramName)

	if param != nil {
		return param
	} else {
		param = &Param{
			Name: paramName,
		}
		p.List = append(p.List, param)
		return param
	}
}

func (p *Params) SetInputParam(paramName string, paramValue any) *Params {

	param := p.FindParamByName(paramName)

	if param != nil {
		param.Value.Value = paramValue
	} else {
		param = &Param{
			Name:      paramName,
			Value:     &variant.Variant{Value: paramValue},
			ParamType: IN,
		}
		p.List = append(p.List, param)
	}

	return p
}

func (p *Params) SetInputParamClob(paramName string, paramValue string) *Params {
	if p.Owner.Connection.Dialect == dbconnect.ORACLE {
		p.SetInputParam(paramName, goOra.Clob{String: paramValue, Valid: StrNotEmpty(paramValue)})
	} else if p.Owner.Connection.Dialect == dbconnect.POSTGRESQL {
		p.SetInputParam(paramName, []byte(paramValue))
	} else {
		p.SetInputParam(paramName, paramValue)
	}
	return p
}

func (p *Params) SetInputParamBlob(paramName string, paramValue []byte) *Params {
	if p.Owner.Connection.Dialect == dbconnect.ORACLE {
		p.SetInputParam(paramName, goOra.Blob{Data: paramValue, Valid: len(paramValue) > 0})
	} else if p.Owner.Connection.Dialect == dbconnect.POSTGRESQL {
		p.SetInputParam(paramName, paramValue)
	} else {
		p.SetInputParam(paramName, paramValue)
	}
	return p
}

func (p *Params) SetOutputParam(paramName string, paramValue any) *Params {

	param := p.FindParamByName(paramName)

	if param != nil {
		param.Value.Value = paramValue
	} else {
		switch paramValue.(type) {
		case int, int8, int16, int32, int64:
			value := int64(0)
			param = &Param{
				Name:      paramName,
				Value:     &variant.Variant{Value: &value},
				ParamType: OUT,
			}
		case float32:
			value := float32(0)
			param = &Param{
				Name:      paramName,
				Value:     &variant.Variant{Value: &value},
				ParamType: OUT,
			}
		case float64:
			value := float64(0)
			param = &Param{
				Name:      paramName,
				Value:     &variant.Variant{Value: &value},
				ParamType: OUT,
			}
		case string:
			value := generateString()
			param = &Param{
				Name:      paramName,
				Value:     &variant.Variant{Value: &value},
				ParamType: OUT,
			}
		case bool:
			value := false
			param = &Param{
				Name:      paramName,
				Value:     &variant.Variant{Value: &value},
				ParamType: OUT,
			}
		case time.Time:
			value := time.Time{}
			param = &Param{
				Name:      paramName,
				Value:     &variant.Variant{Value: &value},
				ParamType: OUT,
			}
		default:
			value := float64(0)
			param = &Param{
				Name:      paramName,
				Value:     &variant.Variant{Value: &value},
				ParamType: OUT,
			}
		}

		p.List = append(p.List, param)
	}

	return p
}

func (p *Params) SetOutputParamSlice(params ...ParamOut) *Params {
	for i := 0; i < len(params); i++ {
		p.SetOutputParam(params[i].Name, params[i].Dest)
	}
	return p
}

func (p *Params) PrintParam() {
	for key, value := range p.List {
		fmt.Println("Colum:", key, "Value:", value.AsValue(), "Type:", reflect.TypeOf(value.AsValue()))
	}
}
