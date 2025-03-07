package preprocesssql

type TEscapeData struct {
	Kind TEscapeKind
	Args TStrings
	Name string
	Func TEscapeFunction
}

func NewEscapedData() *TEscapeData {
	return &TEscapeData{}
}
