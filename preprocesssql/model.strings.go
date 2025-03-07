package preprocesssql

type TStrings []string

func (L *TStrings) Clear() {
	*L = TStrings{}
}
func (L *TStrings) AddStrings(value TStrings) {
	*L = append(*L, value...)
}

func (L *TStrings) Add(value string) {
	*L = append(*L, value)
}
