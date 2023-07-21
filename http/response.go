package http

import "fmt"

type Response struct {
	StatusCode    int
	StatusMessage string
	Body          []byte
	Header        map[string][]string
}

func (R *Response) GetStatusCodeStr() string {
	return fmt.Sprintf("%d", R.StatusCode)
}
func (R *Response) GetStatusMessage() string {
	return R.StatusMessage
}
