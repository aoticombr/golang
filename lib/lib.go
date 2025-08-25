package lib

import (
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	Space     = " "
	EmptyStr  = ""
	LineBreak = "\n"
)

func IsPointer(value interface{}) bool {
	t := reflect.TypeOf(value)
	return t.Kind() == reflect.Ptr
}

/*
CloneBytes:
copia o byte para evitar referencia cruzada de objetos
*/
func CloneBytes(src []byte) []byte {
	clone := make([]byte, len(src))
	copy(clone, src)
	return clone
}

/*
IfThen:
If then em uma unica funcao
*/
func IfThen[T any](condition bool, v1 T, v2 T) T {
	if condition {
		return v1
	} else {
		return v2
	}
}

/*
RemoveFileExtension:
Remove Extensão do arquivo
*/
func RemoveFileExtension(name string) string {
	dotIndex := strings.LastIndex(name, ".")
	if dotIndex == -1 {
		// Caso não encontre um ponto, retorna o nome original
		return name
	}

	return name[:dotIndex]
}

/*
JoinErrors:
Junta uma lista de erros e converte em string
*/
func JoinErrors(errors []error) string {
	var errorMessages []string

	for _, err := range errors {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	return strings.Join(errorMessages, "; ")
}

/*
In:
Verifica se algo esta na lista
*/
func In[T any](valor T, lista []T) bool {
	for _, v := range lista {
		if reflect.DeepEqual(v, valor) {
			return true
		}
	}
	return false
}

/*
InIn:
Procura valores de uma lista dentor de outra lista
Encontrando qualquer valor ele retorna true
*/
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

/*
Guid:
Guid em formato UUID
*/
func Guid() uuid.UUID {
	return uuid.New()
}

/*
GuidString:
Guid em formato string
*/
func GuidString() string {
	return uuid.New().String()
}

/*
IsValidEmail:
Valida se é um email valido
*/
func IsValidEmail(email string) bool {
	// Define a regular expression for validating an email
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ExtractFileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
