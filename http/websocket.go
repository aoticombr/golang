package http

type IWebsocket interface {
	Read(messageType int, body []byte, err error)
}
