package http

type WebSocket struct {
	AutoReconnect    bool
	NumberOfAttempts int
	attempts         int
	connectado       bool
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		AutoReconnect:    true,
		NumberOfAttempts: 10,
		attempts:         0,
		connectado:       false,
	}
}

type IWebsocket interface {
	Read(messageType int, body []byte, err error)
	Error(msg string)
	Msg(msg string)
	Disconect(msg string, limit bool)
}
