package lib

import (
	"fmt"
	"time"
)

/*
StrToDate:
converte string em data dd/mm/yyyy
*/
func StrToDate(str string) (time.Time, error) {
	return time.Parse("02/01/2006", str)
}

/*
DateNilToStr:
Converte data em string dd/mm/yyyy
*/
func DateNilToStr(date *time.Time) string {
	if date == nil {
		return ""
	}
	return date.Format("02/01/2006")
}

/*
FormatDate:
Formata Data
*/
func FormatDate(format string, date time.Time) string {
	return date.Format(format)
}

/*
intFromHex:
int para HEX
*/
func intFromHex(hexStr string) int {
	var result int
	fmt.Sscanf(hexStr, "%x", &result)
	return result
}
