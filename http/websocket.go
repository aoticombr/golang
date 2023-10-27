package http

type IWebsocket interface {
	Read(messageType int, body []byte, err error)
	Error(msg string)
	Msg(msg string)
}
