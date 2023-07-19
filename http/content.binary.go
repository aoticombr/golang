package http

type ContentBinary struct {
	Name     string
	FileName string
	Value    []byte
}
type ListContentBinary []*ContentBinary

func (L *ListContentBinary) Add(name string, fileName string, value []byte) {
	*L = append(*L, &ContentBinary{
		Name:     name,
		FileName: fileName,
		Value:    value,
	})
}
func (L *ListContentBinary) Clear() {
	*L = []*ContentBinary{}
}
func NewListContentBinary() ListContentBinary {
	return make(ListContentBinary, 0)
}
