package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
FormatDatePoint:
Formata ponteiro de data
*/
func FormatDatePoint(format string, date *time.Time) *string {
	if date == nil {
		return nil
	}
	resultado := date.Format(format)
	return &resultado

}

/*
IntToStr:
Converte int para string
*/
func IntToStr(value int) string {
	return strconv.Itoa(value)
}

/*
Int64ToStr:
Converte int64 para string
*/
func Int64ToStr(value int64) string {
	return strconv.FormatInt(value, 10)
}

/*
StrToInt64:
Converte string para int64
*/
func StrToInt64(value string) int64 {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Println("Erro ao converter a string para int64:", err)
		return 0
	}
	return num
}

/*
StrToInt:
Converte string para int
*/
func StrToInt(value string) int {
	num, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("Erro ao converter a string para int:", err)
		return 0
	}
	return num
}

/*
DateToStr:
Converte data em string dd/mm/yyyy
*/
func DateToStr(date time.Time) string {
	return date.Format("02/01/2006")
}

/*
FormatDecimal:
Exp:

	// Convertendo o valor float para uma string formatada
    formattedValue := strconv.FormatFloat(valor, 'f', 2, 64)

    // Substituindo o ponto por vírgula na parte decimal
    formattedValue = formatDecimal(formattedValue)
*/

func FormatDecimal(value string) string {
	// Dividindo o número em parte inteira e parte decimal
	parts := strings.Split(value, ".")
	if len(parts) == 1 {
		// Se não houver parte decimal, retorna apenas a parte inteira
		return parts[0]
	}
	// Formatando a parte decimal com vírgula
	decimalPart := parts[1]
	return parts[0] + "," + decimalPart
}

/*
FormatFloat:
Exp:

	comum.FormatFloat("%%%d.%df", mes.Valor, 4, 2)
	comum.FormatFloat("%%%d.2f", item.Qtde_venda, 14, 2)
	comum.FormatFloat("%%%d.2f", item.Qtde_venda, 14, 2)
*/
func FormatFloat(format string, value float64, width, precision int) string {
	formatString := fmt.Sprintf(format, width, precision)
	return fmt.Sprintf(formatString, value)
}

/*
FloatToStr:
Exp:

	FloatToStr(-4613.70, 2, 64) = -4613.70
*/

func FloatToStr(value float64, width, precision int) string {
	return strconv.FormatFloat(value, 'f', width, precision)
}

/*
FloatToStrPTBR:
Exp:

	FloatToStr(-4613.70, 2, 64) = -4613,70
*/
func FloatToStrPTBR(value float64, width, precision int) string {
	formattedValue1 := FloatToStr(value, width, precision)
	formattedValue2 := FormatDecimal(formattedValue1)
	return formattedValue2
}

/*
FloatToInt:
Exp:

	FloatToInt(-4613.70) = -4613
*/
func FloatToInt(value float64) int {
	return int(value)
}

/*
FormatMatrixFloat:
FormatMatrixFloat(numero, 12, 2) = +000000051.01
FormatMatrixFloat(numero, 15, 3) = +00000000051.010
*/
func FormatMatrixFloat(numero float64, digitos, deci int) string {
	// Converte o número para uma string formatada com 2 casas decimais
	digt := digitos + deci
	if deci > 0 {
		digt += 1
	}
	formatString := fmt.Sprintf("+%%0%d.%df", digt, deci)
	strNumero := fmt.Sprintf(formatString, numero)
	return strNumero
}

// formato 000.000.000-00
func FormatCPF(cpf string) string {
	cpfNumber := OnlyNumber(cpf)

	if len(cpfNumber) == 11 {
		cpf = cpfNumber[0:3] + "." + cpfNumber[3:6] + "." + cpfNumber[6:9] + "-" + cpfNumber[9:11]
	}

	return cpf
}
