package preprocesssql

import "strings"

type TParams struct {
	Items    []*TParam
	BindMode TParamBindMode
	Markers  TStrings
}

func (ps *TParams) Clear() {
	ps.Items = []*TParam{}
}
func (ps *TParams) Add(value *TParam) {
	ps.Items = append(ps.Items, value)
}
func (ps *TParams) Count() int {
	return len(ps.Items)
}
func (ps *TParams) NewParam() *TParam {
	p := NewParam()
	ps.Add(p)
	return p
}
func (ps *TParams) FindParam(name string) *TParam {
	for _, param := range ps.Items {
		if strings.EqualFold(param.Name, name) {
			return param
		}
	}
	return nil
}
func (ps *TParams) IndexOf(name string) int {
	for i, param := range ps.Items {
		if strings.EqualFold(param.Name, name) {
			return i
		}
	}
	return -1
}

func NewParams() *TParams {
	return &TParams{}
}
