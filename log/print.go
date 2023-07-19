package log

import (
	"fmt"
	"log"
	"runtime"
)

var color = true

func init() {
	if runtime.GOOS == "windows" {
		color = false
	}
}

const (
	resetColor = "\033[0m"
	ErrorColor = "\033[31m"
	WarnColor  = "\033[33m"
)

func PrintlnInfo(logger *log.Logger, v interface{}) {
	logger.Println(fmt.Sprint(v))
}

func PrintlnWarn(logger *log.Logger, v interface{}) {
	if color {
		logger.Println(WarnColor + fmt.Sprint(v) + resetColor)
		return
	}
	logger.Println(fmt.Sprint(v))
}

func PrintlnError(logger *log.Logger, v interface{}) {
	if color {
		logger.Println(ErrorColor + fmt.Sprint(v) + resetColor)
		return
	}
	logger.Println(fmt.Sprint(v))
}
