package http

import "fmt"

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

func GetStrFromMethod(methodStr string) (TMethod, error) {
	switch methodStr {
	case "GET":
		return M_GET, nil
	case "POST":
		return M_POST, nil
	case "PUT":
		return M_PUT, nil
	case "DELETE":
		return M_DELETE, nil
	case "HEAD":
		return M_HEAD, nil
	case "OPTIONS":
		return M_OPTIONS, nil
	case "TRACE":
		return M_TRACE, nil
	case "PATCH":
		return M_PATCH, nil
	default:
		return -1, fmt.Errorf("Método HTTP não reconhecido!")
	}
}

type TTransport = int

const (
	TNenhum TTransport = 0
	TSSL    TTransport = 1
	TTLS    TTransport = 2
	TSSLTLS TTransport = 3
)
