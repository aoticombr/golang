package http

import ST "github.com/aoticombr/golang/stringlist"

type ContentText struct {
	Name  string
	Value *ST.Strings
}

type ListContentText []*ContentText

func (L *ListContentText) Add(Name string, value *ST.Strings) {
	*L = append(*L, &ContentText{
		Name:  Name,
		Value: value,
	})
}
func (L *ListContentText) Clear() {
	*L = []*ContentText{}
}
func NewListContentText() ListContentText {
	return make(ListContentText, 0)
}
