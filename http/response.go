package http

type Response struct {
	StatusCode    int
	StatusMessage string
	Body          []byte
	Header        map[string][]string
}
