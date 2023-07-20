package http

type Header struct {
	Accept          string
	AcceptCharset   string
	AcceptEncoding  string
	AcceptLanguage  string
	Authorization   string
	Charset         string
	ContentType     string
	ContentLength   string
	ContentEncoding string
	ContentVersion  string
	ContentLocation string
	ExtraFields     Fields
}

func (H *Header) AddField(fieldName string, fieldValue string) {
	if fieldName == "Accept" {
		H.Accept = fieldValue
	} else if fieldName == "Accept-Charset" {
		H.AcceptCharset = fieldValue
	} else if fieldName == "Accept-Encoding" {
		H.AcceptEncoding = fieldValue
	} else if fieldName == "Accept-Language" {
		H.AcceptLanguage = fieldValue
	} else if fieldName == "Authorization" {
		H.Authorization = fieldValue
	} else if fieldName == "Charset" {
		H.Charset = fieldValue
	} else if fieldName == "Content-Type" {
		H.ContentType = fieldValue
	} else if fieldName == "Content-Length" {
		H.ContentLength = fieldValue
	} else if fieldName == "Content-Encoding" {
		H.ContentEncoding = fieldValue
	} else if fieldName == "Content-Version" {
		H.ContentVersion = fieldValue
	} else if fieldName == "Content-Location" {
		H.ContentLocation = fieldValue
	} else {
		H.ExtraFields.Add(fieldName, fieldValue)
	}
}

func NewHeader() *Header {
	H := &Header{
		Accept:      "*/*",
		ExtraFields: NewFields(),
	}
	return H
}
