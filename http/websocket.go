package http

type Status int

const (
	OPEN Status = iota
	CLOSED
	CONNECTING
	STOP
)

type WebSocket struct {
	AutoReconnect    bool
	NumberOfAttempts int
	attempts         int
	connect          Status
}

func (ws *WebSocket) Connect() Status {
	return ws.connect
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		AutoReconnect:    true,
		NumberOfAttempts: 10,
		attempts:         0,
		connect:          STOP,
	}
}

type IWebsocket interface {
	Read(messageType int, body []byte, err error)
	Error(msg string)
	Msg(msg string)
	Disconect(msg string, limit bool)
}
