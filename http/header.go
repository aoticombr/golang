package http

import "strings"

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
func (H *Header) GetAllFields() map[string]string {
	headerValues := make(map[string]string)

	if H.Accept != "" {
		headerValues["Accept"] = H.Accept
	}
	if H.Accept != "" {
		headerValues["Accept-Charset"] = H.AcceptCharset
	}
	if H.Accept != "" {
		headerValues["Accept-Encoding"] = H.AcceptEncoding
	}
	if H.Accept != "" {
		headerValues["Accept-Language"] = H.AcceptLanguage
	}
	if H.Accept != "" {
		headerValues["Authorization"] = H.Authorization
	}
	if H.Accept != "" {
		headerValues["Charset"] = H.Charset
	}
	if H.Accept != "" {
		headerValues["Content-Type"] = H.ContentType
	}
	if H.Accept != "" {
		headerValues["Content-Length"] = H.ContentLength
	}
	if H.Accept != "" {
		headerValues["Content-Encoding"] = H.ContentEncoding
	}
	if H.Accept != "" {
		headerValues["Content-Version"] = H.ContentVersion
	}
	if H.Accept != "" {
		headerValues["Content-Location"] = H.ContentLocation
	}

	// Adicionando os campos extras do cabeçalho (ExtraFields)
	for fieldName, fieldValues := range H.ExtraFields {
		if len(fieldValues) > 0 {
			// Se houver múltiplos valores para o campo, concatenamos eles em uma única string
			headerValues[fieldName] = strings.Join(fieldValues, "; ")
		}
	}

	return headerValues
}

func NewHeader() *Header {
	H := &Header{
		Accept:      "*/*",
		ExtraFields: NewFields(),
	}
	return H
}
