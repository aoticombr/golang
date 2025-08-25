package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aoticombr/golangprivate/models"
)

/*
FileExists:
Verifica se o arquivo existe
*/
func FileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

/*
ByteToSaveFile:
Converte Byte para Arquivo
*/
func ByteToSaveFile(nomeArquivo string, dados []byte) error {
	// Cria ou abre o arquivo para escrita
	arquivo, err := os.Create(nomeArquivo)
	if err != nil {
		return fmt.Errorf("erro ao Criar o Arquivo %w", err)
	}
	defer arquivo.Close()

	// Escreve os dados no arquivo
	_, err = arquivo.Write(dados)
	if err != nil {
		return fmt.Errorf("erro ao Escrever o Arquivo %w", err)
	}

	return nil
}

/*
ListarArquivosPastaTodos:
Lista arquivos de uma pasta informada
*/
func ListarArquivosPastaTodos(pasta string) ([]MFile, error) {
	var arquivos []MFile

	err := filepath.Walk(pasta, func(caminho string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar diretórios
		if info.IsDir() {
			return nil
		}

		// Criar uma instância da struct Arquivo com informações do arquivo
		arquivo := MFile{
			Nome:            info.Name(),
			CaminhoCompleto: caminho,
			Tamanho:         info.Size(),
		}

		// Adicionar à lista de arquivos
		arquivos = append(arquivos, arquivo)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return arquivos, nil
}

/*
ListarArquivosPastaFiltro:
Lista arquivos de uma pasta, porem voce pode filtrar por tipo de arquivos e tambem colocar um prefixo de procura
exemplo: ListarArquivosPastaFiltro("c:\", "VRIN*", ".txt")
*/
func ListarArquivosPastaFiltro(pasta, prefixo, extensao string) ([]models.File, error) {
	var arquivos []models.File

	err := filepath.Walk(pasta,
		func(caminho string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Verifica se é um arquivo e se o nome começa com o prefixo e tem a extensão desejada
			if !info.IsDir() && strings.HasPrefix(info.Name(), prefixo) && strings.HasSuffix(info.Name(), extensao) {
				arquivo := models.File{
					Nome:            info.Name(),
					CaminhoCompleto: caminho,
					Tamanho:         info.Size(),
				}

				// Adicionar à lista de arquivos
				arquivos = append(arquivos, arquivo)
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return arquivos, nil
}
