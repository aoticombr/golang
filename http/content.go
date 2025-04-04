package http

import (
	"strings"
)

type EncType int

const (
	ET_NONE                  EncType = 0
	ET_FORM_DATA             EncType = 1
	ET_X_WWW_FORM_URLENCODED EncType = 2
	ET_RAW                   EncType = 3
	ET_BINARY                EncType = 4
	ET_GRAPHQL               EncType = 5
	ET_WEB_SERVICE           EncType = 6
)

type ContentType int

const (
	CT_NONE                  ContentType = 0
	CT_TEXT                  ContentType = 1
	CT_JAVASCRIPT            ContentType = 2
	CT_JSON                  ContentType = 3
	CT_HTML                  ContentType = 4
	CT_XML                   ContentType = 5
	CT_MULTIPART_FORM_DATA   ContentType = 6
	CT_X_WWW_FORM_URLENCODED ContentType = 7
	CT_BINARY                ContentType = 8
	CT_SOAPXML               ContentType = 9
	CT_PDF                   ContentType = 10
)

func GeContentTypeStr(value ContentType) string {

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
	case CT_SOAPXML:
		return "application/soap+xml"
	case CT_PDF:
		return "application/pdf"
	default:
		return ""
	}
}

func GetContentTypeFromString(str string) ContentType {
	//fmt.Println("GetContentTypeFromString: '%s'", strings.ToLower(str))
	switch strings.ToLower(str) {
	case "text/plain":
		return CT_TEXT
	case "application/javascript":
		return CT_JAVASCRIPT
	case "application/json":
		return CT_JSON
	case "text/html":
		return CT_HTML
	case "application/xml":
		return CT_XML
	case "multipart/form-data":
		return CT_MULTIPART_FORM_DATA
	case "application/x-www-form-urlencoded":
		return CT_X_WWW_FORM_URLENCODED
	case "application/octet-stream":
		return CT_BINARY
	default:
		return CT_NONE
	}
}

type Content interface {
}
