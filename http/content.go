package http

import ST "github.com/aoticombr/golang/component"

type TContentType int

const (
	CT_NENHUM                TContentType = 0
	CT_TEXT                  TContentType = 1
	CT_JAVASCRIPT            TContentType = 2
	CT_JSON                  TContentType = 3
	CT_HTML                  TContentType = 4
	CT_XML                   TContentType = 5
	CT_MULTIPART_FORM_DATA   TContentType = 6
	CT_X_WWW_FORM_URLENCODED TContentType = 7
	CT_BINARY                TContentType = 8
)

func GetContentTypeStr(value TContentType) string {

	switch value {
	case CT_TEXT:
		return "text/plain"
	case CT_JAVASCRIPT:
		return "application/javascript"
	case CT_JSON:
		return "application/json"
	case CT_HTML:
		return "text/html"
	case CT_XML:
		return "application/xml"
	case CT_MULTIPART_FORM_DATA:
		return "multipar/form-data"
	case CT_X_WWW_FORM_URLENCODED:
		return "application/x-www-form-urlencoded"
	case CT_BINARY:
		return "application/octet-stream"
	default:
		return ""
	}
}

type Content interface {
}
type ContentFormField struct {
	FieldName  string
	FieldValue string
}
type ContentFile struct {
	FieldName   string
	Content     []byte
	ContentType string
}
type ContentText struct {
	Value ST.Strings
}
type ContentBinary struct {
	Value []byte
}
type ListContentBinary []ContentBinary
type ListContentText []ContentText
type ListContentFile []ContentFile
type ListContentFormField []ContentFormField
