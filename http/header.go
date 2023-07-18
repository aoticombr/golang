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

	ItensFormField   ListContentFormField
	ItensSubmitFile  ListContentFile
	ItensContentText ListContentText
	ItensContentBin  ListContentBinary
}
