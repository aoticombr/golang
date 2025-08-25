package lib

import "strings"

/*
Copy:
Imita a funcao Copy delphi
*/
func Copy(value string, col_ini, qtde_caracteres int) string {
	t_value := len(value)
	if t_value > 0 { //precisa
		if col_ini < 1 {
			col_ini = 1
		}
		if col_ini < t_value {
			t := t_value
			if col_ini > 1 {
				t = t_value - col_ini + 1
			}
			if t >= qtde_caracteres {
				return value[col_ini-1 : col_ini-1+qtde_caracteres]
			} else {
				return value[col_ini-1 : col_ini-1+t]
			}
		}
	}
	return ""
}

/*
Trim:
Remove espaços
*/
func Trim(value string) string {
	return strings.TrimSpace(value)
}

func Exclude[T comparable](slice []T, value T) []T {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// Include adds a given value to a slice, if it doesn't already exist.
func Include[T comparable](slice []T, value T) []T {
	for _, v := range slice {
		if v == value {
			return slice
		}
	}
	return append(slice, value)
}

/*
Pos:
Esta função Pos aceita dois argumentos: a substring que você está procurando e a
string na qual você está procurando. Ela retorna o índice baseado em 1 da primeira
ocorrência da substring na string. Se a substring não for encontrada, ela retorna 0
Comentario:

	por mais que exista strings.Contains essa aqui retorna posicao onde encontrou o dado

caso de uso:

	use para recortar um dado apos encontrar um dados desejado
	encontrar marcadores como <> ,</> para recortar dados
*/
func Pos(substring, s string) int {
	index := strings.Index(s, substring)
	if index == -1 {
		return 0 // Se a substring não for encontrada, retorna 0
	}
	return index + 1 // Se a substring for encontrada, retorna o índice baseado em 1
}

/*
QuotedStr:
coloca aspas simples entre uma string
*/
func QuotedStr(s string) string {
	// Replace single quotes with two single quotes
	s = strings.ReplaceAll(s, "'", "''")
	// Wrap the string with single quotes
	return "'" + s + "'"
}
func IsQuotedBase(value string, ACh1, ACh2 byte) bool {
	return (len(value) > 2) && In(ACh1, []byte{0, ' '}) && !In(ACh2, []byte{0, ' '}) && (value[1] == ACh1) && (value[len(value)] == ACh2)
}

func UnQuoteBase(value string, ACh1, ACh2 byte) string {
	if IsQuotedBase(value, ACh1, ACh2) {
		return Copy(value, 2, len(value)-2)
	} else {
		return value
	}
}

func AnsiQuotedStr(S string, Quote byte) string {
	return string(Quote) + S + string(Quote)
}

func AnsiCompareText(str1, str2 string) int {
	// Convertendo ambas as strings para minúsculas
	str1Lower := strings.ToLower(str1)
	str2Lower := strings.ToLower(str2)

	// Comparando as strings minúsculas
	if str1Lower == str2Lower {
		return 0
	} else if str1Lower < str2Lower {
		return -1
	} else {
		return 1
	}
}
