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

func (H *Request) CopyBody(value []byte) {
	H.Body = make([]byte, len(value))
	copy(H.Body, value)
}

func (H *Request) AddFormField(fieldName string, fieldValue string) {
	H.ItensFormField.Add(fieldName, fieldValue)
}
func (H *Request) AddSubmitFile(key string, filename string, contentType string, content []byte) {
	H.ItensSubmitFile.Add(key, filename, contentType, content)
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
