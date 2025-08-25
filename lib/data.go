package lib

import "time"

const Format1 = "2006-01-02T15:04:05"
const Format2 = "2006-01-02T15:04:05.000" //2025-02-30T00:00:00.000

/*
EndDayOfYearMonth:
Ultimo dia do Ano/Mes , exemp: 2024/12 = dia é 31
*/
func EndDayOfYearMonth(year int, month time.Month) time.Time {
	// Criar o primeiro dia do próximo mês
	firstDayOfNextMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC)

	// Subtrair um dia para obter o último dia do mês atual
	lastDayOfCurrentMonth := firstDayOfNextMonth.Add(-time.Second)

	return lastDayOfCurrentMonth
}

func StrIso8601ToDate1(str string) time.Time {
	t, _ := time.Parse(Format1, str)
	return t
}
func StrIso8601ToDate2(str string) time.Time {
	t, _ := time.Parse(Format2, str)
	return t
}
