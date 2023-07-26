package jsonconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Schema struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Schema string `json:"schema"`
	SID    string `json:"sid"`
	Ativo  bool   `json:"ativo"`
}

type Boot struct {
	Name    string   `json:"name"`
	Schemas []Schema `json:"schemas"`
	Ativo   bool     `json:"ativo"`
}

type Gateway struct {
	Protocolo string `json:"protocolo"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Ativo     bool   `json:"ativo"`
}

type Api struct {
	Name      string   `json:"name"`
	Protocolo string   `json:"protocolo"`
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	Gateway   Gateway  `json:"gateway"`
	Schemas   []Schema `json:"schemas"`
	Ativo     bool     `json:"ativo"`
}

type Config struct {
	Boots []Boot `json:"boots"`
	Apis  []Api  `json:"apis"`
	Path  string `json:"path"`
}

func main() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Erro ao fazer o parsing do JSON:", err)
		return
	}

	// Agora você tem os dados do JSON armazenados na variável 'config'
	// Você pode acessá-los conforme necessário.

	// Exemplo de como acessar o nome do primeiro serviço de boot:
	if len(config.Boots) > 0 {
		firstBootName := config.Boots[0].Name
		fmt.Println("Nome do primeiro serviço de boot:", firstBootName)
	}

	// Exemplo de como acessar o nome do primeiro plugin de API:
	if len(config.Apis) > 0 {
		firstApiName := config.Apis[0].Name
		fmt.Println("Nome do primeiro plugin de API:", firstApiName)
	}

	// Exemplo de como acessar a senha do primeiro schema do primeiro serviço de boot:
	if len(config.Boots) > 0 && len(config.Boots[0].Schemas) > 0 {
		firstBootFirstSchemaPass := config.Boots[0].Schemas[0].Pass
		fmt.Println("Senha do primeiro schema do primeiro serviço de boot:", firstBootFirstSchemaPass)
	}

	// E assim por diante, você pode acessar os outros campos da estrutura 'Config' conforme necessário.

}
