package preprocesssql

type TParam struct {
	Index           int
	Name            string
	Position        int
	IsCaseSensitive bool
	ParamType       TParamType
}

func NewParam() *TParam {
	return &TParam{}
}
