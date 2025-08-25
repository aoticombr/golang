package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadJsonFileToStruct(PahtFile string, value any) error {
	// Lê o arquivo JSON
	file, err := os.Open(PahtFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Lê o conteudo do arquivo
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Descompacta(Unmarshal) os dados do JSON no objeto ModelDesign
	err = json.Unmarshal(content, value)
	if err != nil {
		return err
	}
	return nil
}
