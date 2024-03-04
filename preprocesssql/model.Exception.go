package preprocesssql

import "fmt"

type FDException struct {
	Msg string
}

func NewFDException(msg string) FDException {
	return FDException{Msg: msg}
}

func (e FDException) Error() string {
	return fmt.Sprintf("FD Exception: %s", e.Msg)
}
