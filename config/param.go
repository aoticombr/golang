package config

import "strings"

type Param struct {
	Marcador string
	Name     string
	Value    string
}

type Params []*Param

func (p *Params) GetParam(Marcador, Name string) (bool, string) {
	for _, v := range *p {
		if strings.ToUpper(v.Marcador) == strings.ToUpper(Marcador) {
			if strings.ToUpper(v.Name) == strings.ToUpper(Name) {
				return true, v.Value
			}
		}
	}
	return false, ""
}
