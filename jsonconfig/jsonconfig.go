package jsonconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JsonConfig struct {
	Name   string
	config *Config
}

func (j *JsonConfig) GetConfig() *Config {
	return j.config
}

func (j *JsonConfig) Load() error {
	data, err := ioutil.ReadFile(j.Name)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return err
	}
	err = json.Unmarshal(data, &j.config)
	if err != nil {
		return fmt.Errorf("Erro ao fazer o parsing do JSON: %s", err)
	}

	// Agora você tem os dados do JSON armazenados na variável 'config'
	// Você pode acessá-los conforme necessário.

	// Exemplo de como acessar o nome do primeiro serviço de boot:
	if len(j.config.Boots) > 0 {
		firstBootName := j.config.Boots[0].Name
		fmt.Println("Nome do primeiro serviço de boot:", firstBootName)
	}

	// Exemplo de como acessar o nome do primeiro plugin de API:
	if len(j.config.Apis) > 0 {
		firstApiName := j.config.Apis[0].Name
		fmt.Println("Nome do primeiro plugin de API:", firstApiName)
	}

	// Exemplo de como acessar a senha do primeiro schema do primeiro serviço de boot:
	if len(j.config.Boots) > 0 && len(j.config.Boots[0].Schemas) > 0 {
		firstBootFirstSchemaPass := j.config.Boots[0].Schemas[0].Pass
		fmt.Println("Senha do primeiro schema do primeiro serviço de boot:", firstBootFirstSchemaPass)
	}
	return nil
	// E assim por diante, você pode acessar os outros campos da estrutura 'Config' conforme necessário.

}

func NewJsonConfig() *JsonConfig {
	return &JsonConfig{}
}
