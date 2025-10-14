package lib

import (
	"os"
	"strconv"
	"strings"
)

/*
GetEnvInt:
obte int de .env
*/
func GetEnvInt(value string) int {
	str := os.Getenv(value)
	intValue, _ := strconv.Atoi(str)
	return intValue
}

/*
GetEnvString:
obte string de .env
*/
func GetEnvString(value string) string {
	return os.Getenv(value)
}

/*
GetEnvBool:
obte boleano de .env
*/
func GetEnvBool(value string) bool {
	str := os.Getenv(value)
	return strings.ToUpper(str) == "TRUE"
}
