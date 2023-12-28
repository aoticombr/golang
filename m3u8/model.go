package m3u8

import (
	"fmt"
	"strconv"
	"strings"
)

type Resolution struct {
	Width  int
	Height int
}
type Attribute struct {
	Resolution
	Uri string
}
type File struct {
	Kind string
	Name string
}

func NewAttribute(value string) (*Attribute, error) {

	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("formato de entrada inválido")
	}

	// Extrair resolução
	resolutionParts := strings.Split(parts[0], "x")
	if len(resolutionParts) != 2 {
		return nil, fmt.Errorf("formato de resolução inválido")
	}

	width, err := strconv.Atoi(resolutionParts[0])
	if err != nil {
		return nil, fmt.Errorf("erro ao converter largura para inteiro: %v", err)
	}

	height, err := strconv.Atoi(resolutionParts[1])
	if err != nil {
		return nil, fmt.Errorf("erro ao converter altura para inteiro: %v", err)
	}

	// Criar instância da struct Attribute
	attribute := &Attribute{
		Resolution: Resolution{
			Width:  width,
			Height: height,
		},
		Uri: value,
	}

	return attribute, nil

}
