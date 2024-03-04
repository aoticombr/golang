package preprocesssql

import (
	"reflect"
	"strconv"
	"strings"
)

// TADCharSet é um tipo para representar um conjunto de caracteres

// ADInSet verifica se um caractere está presente em um conjunto de caracteres
func ADInSet(AChar byte, ASet TADCharSet) bool {
	_, ok := ASet[AChar]
	return ok
}

func In[T any](valor T, lista []T) bool {
	for _, v := range lista {
		if reflect.DeepEqual(v, valor) {
			return true
		}
	}
	return false
}
func InIn[T any](lista1 []T, lista2 []T) bool {
	for _, a := range lista1 {
		for _, b := range lista2 {
			if reflect.DeepEqual(b, a) {
				return true
			}
		}
	}
	return false
}

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

func Move(source []byte, dest []byte, count int) {
	if &source[0] == &dest[0] {
		return // Source = Dest
	}

	// Perform bounds check to prevent index out of range
	if count <= 0 || count > len(source) || count > len(dest) {
		return
	}

	if count <= 8 {
		// Tiny Move (0..8 Byte Move)
		for i := 0; i < count; i++ {
			dest[i] = source[i]
		}
	} else {
		// Large Move
		copy(dest, source)

		// Or use the built-in copy function:
		// copy(dest[:count], source[:count])
	}
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
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
func QuotedStr(s string) string {
	// Replace single quotes with two single quotes
	s = strings.ReplaceAll(s, "'", "''")
	// Wrap the string with single quotes
	return "'" + s + "'"
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
