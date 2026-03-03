package http

type Params map[string]string

func (P Params) Add(key string, value string) {
	P[key] = value
}
func (P Params) Set(key string, value string) {
	P[key] = value
}
func (P Params) Get(key string) string {
	return P[key]
}
func (P Params) Clear() {
	P = make(map[string]string)
}
func NewParams() Params {
	return make(Params)
}
