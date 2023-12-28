package m3u8

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

func extractValuesFromBody(body string) []string {
	var values []string

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		// Remove espaços em branco no início e no final da linha
		line = strings.TrimSpace(line)

		// Ignorar linhas vazias ou que começam com #
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Adicionar o valor ao array
		values = append(values, line)
	}

	return values
}

func parseAttribute(input string) (*Attribute, error) {
	parts := strings.Split(input, "/")
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
		Uri: parts[2],
	}

	return attribute, nil
}

func replacePathInURL(baseURL, newPath string) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Obter o caminho original
	basePath := parsedURL.Path

	// Remover a última parte do caminho (se existir)
	basePath = strings.TrimSuffix(basePath, "/playlist.m3u8")

	// Concatenar o novo caminho ao caminho original
	newPath = path.Join(basePath, newPath)

	// Substituir o caminho na URL
	parsedURL.Path = newPath

	return parsedURL.String(), nil
}

func GetFile(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if resp.StatusCode == 403 {
		return nil
	}
	//fmt.Println(resp.Status)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Open the file in append mode, or create it if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Append the new data to the file
	_, err = file.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func GetFileByte(url string, old []byte, headers map[string]string) ([]byte, error) {
	// Criar um cliente HTTP personalizado com cabeçalhos adicionais
	client := &http.Client{}

	// Criar uma solicitação HTTP personalizada com o método GET e o URL fornecido
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Adicionar cabeçalhos à solicitação
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Executar a solicitação HTTP personalizada
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Concatenar os dados existentes com os novos dados
	newData := append(old, body...)

	return newData, nil
}
func DownloadByte(List []string) ([]byte, error) {
	var bt []byte
	headers := map[string]string{
		"Referer": "https://google.com.br",
	}

	for _, file := range List {
		btnew, err := GetFileByte(file, bt, headers)
		if err != nil {
			return nil, err
		}
		bt = btnew
	}
	return bt, nil
}

func Download(List []string, File *File) {
	for _, file := range List {
		err := GetFile(file, File.Name)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SaveByteToFile(fileName string, data []byte) error {
	err := ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
