package http

import ST "github.com/aoticombr/golang/component"

type Request struct {
	Header           *Header
	ItensFormField   ListContentFormField
	ItensSubmitFile  ListContentFile
	ItensContentText ListContentText
	ItensContentBin  ListContentBinary
	Body             []byte
}

func (H *Request) AddFormField(fieldName string, fieldValue string) {
	H.ItensFormField.Add(fieldName, fieldValue)
}
func (H *Request) AddSubmitFile(fieldName string, contentType string, content []byte) {
	H.ItensSubmitFile.Add(fieldName, contentType, content)
}
func (H *Request) AddContentText(Name string, value *ST.Strings) {
	H.ItensContentText.Add(Name, value)
}
func (H *Request) AddContentBin(name string, filename string, value []byte) {
	H.ItensContentBin.Add(name, filename, value)
}

func NewRequest() *Request {
	R := &Request{
		Header:           NewHeader(),
		ItensFormField:   NewListContentFormField(),
		ItensSubmitFile:  NewListContentFile(),
		ItensContentText: NewListContentText(),
		ItensContentBin:  NewListContentBinary(),
	}
	return R
}
