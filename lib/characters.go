package lib

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

/*
RandomString:
Gera uma string do tamanho desejado com caracteres aleatorios, otimo para um gerador de senhas
*/

func RandomString(qtde int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, qtde)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)

}

/*
LimitStr:
Corta a string caso exceda
*/
func LimitStr(value string, limit int) string {
	if len(value) > limit {
		return value[:limit]
	}
	return value
}

/*
isANSI:
Verifica seu o formato é ANSI
*/
func isANSI(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

/*
OnlyAlfaNumber:
Deixa apenas Number/Alfa passar
*/
func OnlyAlfaNumber(s string) string {
	//fmt.Println(s)
	result := []rune(s)

	// Cria uma nova fatia para armazenar os caracteres desejados
	var filtered []rune

	for i := len(result) - 1; i >= 0; i-- {
		if unicode.IsLetter(result[i]) || unicode.IsNumber(result[i]) {
			// Adiciona os caracteres desejados à nova fatia
			filtered = append([]rune{result[i]}, filtered...)
		}
	}

	// Atribui a nova fatia à variável original
	result = filtered

	return string(result)
}

/*
RemoveAcentos:
Remove acentos
*/
func RemoveAcentos(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		return s
	}
	return output
}

/*
RemoveCaracteres:
Remove tudo que não é letra, numero, ou espaco, pode remover acentos tambem caso[pRemoveAcentos]=(true)
*/
func RemoveCaracteres(texto string, pRemoveAcentos bool) string {
	resultado := texto
	if pRemoveAcentos {
		resultado = RemoveAcentos(texto)
	}

	for i := len(resultado) - 1; i >= 0; i-- {
		if !unicode.IsDigit(rune(resultado[i])) && !unicode.IsLetter(rune(resultado[i])) && resultado[i] != ' ' {
			resultado = resultado[:i] + resultado[i+1:]

		}
	}

	return resultado
}

/*
Lpad:
Completa Caracteres a Esquerda(Left)
*/
func Lpad(s string, length int, padChar string) string {
	return strings.Repeat(padChar, length-len(s)) + s
}

/*
Rpad:
Completa Caracteres a Direita(Right)
*/
func Rpad(s string, length int, padChar string) string {
	return s + strings.Repeat(padChar, length-len(s))
}

/*
RemoverQuebraLinha:
Função para remover quebras de linha \r ou \n
*/
func RemoverQuebraLinha(texto string) string {
	texto = strings.ReplaceAll(texto, "\r", "") // Remove \r
	texto = strings.ReplaceAll(texto, "\n", "") // Remove \n
	return texto
}

/*
Espacos:
Retornar string cheia de espaços
*/
func Espacos(quantidade int) string {
	return strings.Repeat(" ", quantidade)
}

/*
RepetirString:
Repete uma string x vezes
*/
func RepetirString(c string, quantidade int) string {
	return strings.Repeat(c, quantidade)
}

/*
PreencherEsquerda:
Preenche a esqueda com o caracter informado, caso exceda o tamanho informado ele recorta o valor
*/
func PreencherEsquerda(s string, tamanho int, caracter string) string {
	l := len(s)
	if l > tamanho {
		return s[:tamanho] //recorda caso exceda
	} else if l < tamanho {
		return s + RepetirString(caracter, tamanho-l)
	}
	return s
}

/*
PreencherDireita:
Preenche a direita com o caracter informado, caso exceda o tamanho informado ele recorta o valor
*/
func PreencherDireita(s string, tamanho int, caracter string) string {
	l := len(s)
	if l > tamanho {
		return s[:tamanho] //recorda caso exceda
	} else if l < tamanho {
		return RepetirString(caracter, tamanho-l) + s
	}
	return s
}

/*
OnlyNumber:
Deixa apenas numeros passar
*/
func OnlyNumber(s string) string {
	result := []rune(s)

	// Cria uma nova fatia para armazenar os caracteres desejados
	var filtered []rune

	for i := len(result) - 1; i >= 0; i-- {
		if unicode.IsNumber(result[i]) {
			// Adiciona os caracteres desejados à nova fatia
			filtered = append([]rune{result[i]}, filtered...)
		}
	}

	// Atribui a nova fatia à variável original
	result = filtered

	return string(result)
}

/*
OnlyAlfa:
Deixa apenas Alfa passar
*/
func OnlyAlfa(s string) string {
	result := []rune(s)

	// Cria uma nova fatia para armazenar os caracteres desejados
	var filtered []rune

	for i := len(result) - 1; i >= 0; i-- {
		if unicode.IsLetter(result[i]) {
			// Adiciona os caracteres desejados à nova fatia
			filtered = append([]rune{result[i]}, filtered...)
		}
	}

	// Atribui a nova fatia à variável original
	result = filtered

	return string(result)
}

/*
OnlyEmail:
Deixa apenas caracteres de email
*/
func OnlyEmail(email string) string {
	const caracteresValidos = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@.-_"
	email = RemoveAcentos(strings.ToLower(email))
	var result strings.Builder
	for _, c := range email {
		if strings.ContainsRune(caracteresValidos, c) {
			result.WriteRune(c)
		}
	}
	return result.String()
}

func ConvertToStringPointers(input []string) []*string {
	var result []*string
	for i := range input {
		result = append(result, &input[i]) // Usamos o índice para obter um ponteiro para cada elemento original
	}
	return result
}
