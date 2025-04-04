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
	CT_ZIP                   ContentType = 11
	CT_PNG                   ContentType = 12
	CT_JPEG                  ContentType = 13
	CT_GIF                   ContentType = 14
	CT_SVGXML                ContentType = 15
	CT_MPEG                  ContentType = 16
	CT_OGG                   ContentType = 17
	CT_MP4                   ContentType = 18
	CT_WEBM                  ContentType = 19
	CT_CSS                   ContentType = 20
	CT_CSV                   ContentType = 21
	CT_TSV                   ContentType = 22
	CT_YAML                  ContentType = 23
	CT_TAR                   ContentType = 24
	CT_RTF                   ContentType = 25
	CT_WAV                   ContentType = 26
	CT_FLAC                  ContentType = 27
	CT_AVI                   ContentType = 28
	CT_MOV                   ContentType = 29
	CT_BMP                   ContentType = 30
	CT_WEBP                  ContentType = 31
	CT_TIFF                  ContentType = 32
	CT_EOT                   ContentType = 33
	CT_TTF                   ContentType = 34
	CT_WOFF                  ContentType = 35
	CT_WOFF2                 ContentType = 36
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
	case CT_ZIP:
		return "application/zip"
	case CT_PNG:
		return "image/png"
	case CT_JPEG:
		return "image/jpeg"
	case CT_GIF:
		return "image/gif"
	case CT_SVGXML:
		return "image/svg+xml"
	case CT_MPEG:
		return "audio/mpeg"
	case CT_OGG:
		return "audio/ogg"
	case CT_MP4:
		return "video/mp4"
	case CT_WEBM:
		return "video/webm"
	case CT_CSS:
		return "text/css"
	case CT_CSV:
		return "text/csv"
	case CT_TSV:
		return "text/tab-separated-values"
	case CT_YAML:
		return "application/x-yaml"
	case CT_TAR:
		return "application/x-tar"
	case CT_RTF:
		return "application/rtf"
	case CT_WAV:
		return "audio/wav"
	case CT_FLAC:
		return "audio/flac"
	case CT_AVI:
		return "video/x-msvideo"
	case CT_MOV:
		return "video/quicktime"
	case CT_BMP:
		return "image/bmp"
	case CT_WEBP:
		return "image/webp"
	case CT_TIFF:
		return "image/tiff"
	case CT_EOT:
		return "application/vnd.ms-fontobject"
	case CT_TTF:
		return "font/ttf"
	case CT_WOFF:
		return "font/woff"
	case CT_WOFF2:
		return "font/woff2"
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
	case "application/soap+xml":
		return CT_SOAPXML
	case "application/pdf":
		return CT_PDF
	case "application/zip":
		return CT_ZIP
	case "image/png":
		return CT_PNG
	case "image/jpeg":
		return CT_JPEG
	case "image/gif":
		return CT_GIF
	case "image/svg+xml":
		return CT_SVGXML
	case "audio/mpeg":
		return CT_MPEG
	case "audio/ogg":
		return CT_OGG
	case "video/mp4":
		return CT_MP4
	case "video/webm":
		return CT_WEBM
	case "text/css":
		return CT_CSS
	case "text/csv":
		return CT_CSV
	case "text/tab-separated-values":
		return CT_TSV
	case "application/x-yaml":
		return CT_YAML
	case "application/x-tar":
		return CT_TAR
	case "application/rtf":
		return CT_RTF
	case "audio/wav":
		return CT_WAV
	case "audio/flac":
		return CT_FLAC
	case "video/x-msvideo":
		return CT_AVI
	case "video/quicktime":
		return CT_MOV
	case "image/bmp":
		return CT_BMP
	case "image/webp":
		return CT_WEBP
	case "image/tiff":
		return CT_TIFF
	case "application/vnd.ms-fontobject":
		return CT_EOT
	case "font/ttf":
		return CT_TTF
	case "font/woff":
		return CT_WOFF
	case "font/woff2":
		return CT_WOFF2

	default:
		return CT_NONE
	}
}

type Content interface {
}
