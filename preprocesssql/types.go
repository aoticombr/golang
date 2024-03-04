package preprocesssql

type TPreprocessorInstr int
type TNameQuoteLevel int
type TNameQuoteSide int
type TTextEndOfLine int
type TDBMSKind int
type TEscapeKind int
type TCommandKind int
type TParamBindMode int
type TParamType int
type TParamMark int
type TEscapeFunction int
type TEncodeOption int
type TNamePart int

type TListByte []byte
type NameQuoteSides []TNameQuoteSide
type TTextEndOfLines []TTextEndOfLine
type TADCharSet map[byte]bool
type TDBMSKinds []TDBMSKind
type TParamTypes []TParamType
type TNameParts []TNamePart
type TEncodeOptions []TEncodeOption
type TPreprocessorInstrs []TPreprocessorInstr
type TNameQuoteLevels []TNameQuoteLevel

func (lb *TListByte) Clear() {
	*lb = TListByte{}
}

func (L *TPreprocessorInstrs) Clear() {
	*L = TPreprocessorInstrs{}
}
func (L *TPreprocessorInstrs) Add(value TPreprocessorInstr) {
	*L = append(*L, value)
}
func (L *TPreprocessorInstrs) Remove(value TPreprocessorInstr) {
	for i, v := range *L {
		if v == value {
			*L = append((*L)[:i], (*L)[i+1:]...)
			break
		}
	}
}
func (L *TPreprocessorInstrs) Removes(value TPreprocessorInstrs) {
	for _, v := range value {
		L.Remove(v)
	}

}

func NewTNameQuoteLevels() TNameQuoteLevels {
	return TNameQuoteLevels{}
}
