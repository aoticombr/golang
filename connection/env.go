package connection

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvInt(value string) int {
	str := os.Getenv(value)
	intValue, _ := strconv.Atoi(str)
	return intValue
}

func GetEnvString(value string) string {
	return os.Getenv(value)
}

func GetEnvBool(value string) bool {
	str := os.Getenv(value)
	return strings.ToUpper(str) == "TRUE"
}
