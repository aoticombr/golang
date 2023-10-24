package http

type IWebsocket interface {
	read(messageType int, body []byte, err error)
}
