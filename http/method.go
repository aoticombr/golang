package http

type TMethod int

const (
	M_GET     TMethod = 0
	M_POST    TMethod = 1
	M_PUT     TMethod = 2
	M_DELETE  TMethod = 3
	M_HEAD    TMethod = 4
	M_OPTIONS TMethod = 5
	M_TRACE   TMethod = 6
	M_PATCH   TMethod = 7
)

func GetMethodStr(value TMethod) string {
	switch value {
	case M_GET:
		return "GET"
	case M_POST:
		return "POST"
	case M_PUT:
		return "PUT"
	case M_DELETE:
		return "DELETE"
	case M_HEAD:
		return "HEAD"
	case M_OPTIONS:
		return "OPTIONS"
	case M_TRACE:
		return "TRACE"
	case M_PATCH:
		return "PATCH"
	default:
		return ""

	}
}
