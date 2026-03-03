package lib

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

/*
FileToGzipFile
Comprime um arquivo em disco para um arquivo gzip
*/
func FileToGzipFile(nomeArquivo, nomedozip string) error {
	/*#######################################
	     Ler o Arquivo
	#######################################*/
	file, err := os.Open(nomeArquivo)
	if err != nil {
		return err
	}
	defer file.Close()

	/*#######################################
	     Gerar o gzip
	#######################################*/
	outputWriterGZ, err := os.Create(nomedozip) //cria o arquivo gz
	if err != nil {
		return err
	}
	defer outputWriterGZ.Close()
	gzipWriter := gzip.NewWriter(outputWriterGZ) //PASSA O ARQUIVO PARA O GZIP
	defer gzipWriter.Close()
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		return err
	}

	return nil
}

// Função para compactar uma string usando GZIP
func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write([]byte(data))
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buf.Bytes(), nil
}

// Função para descompactar um GZIP para string
func GzipDecompress(compressedData []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var decompressed bytes.Buffer
	_, err = io.Copy(&decompressed, reader)
	if err != nil {
		return "", err
	}
	return decompressed.String(), nil
}

/*
ByteToGzipByte
Esta funcao faz:
1 - pega um []byte coloca um nome nele
2 - Comprime ele em .gz e retorna o []byte (dele em gz)
Ele faz tudo isso em memoria sem gerar arquivo em disco
Otimo para enviar via http, ftp, etc
*/
func ByteToGzipByte(data []byte, filename string) ([]byte, error) {
	// Buffer para armazenar o arquivo gzip em memória
	var buf bytes.Buffer

	// Criar um gzip writer
	gz := gzip.NewWriter(&buf)
	defer gz.Close()

	// Configurar o cabeçalho do gzip para incluir o nome do arquivo
	if filename != "" {
		gz.Name = filename
	}

	// Escrever os dados comprimidos no gzip writer
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	// Fechar o gzip writer para finalizar a compressão
	if err := gz.Close(); err != nil {
		return nil, err
	}

	// Retornar os bytes do arquivo gzip
	return buf.Bytes(), nil
}

/*
DecompressGzipToSaveOneFile
Esta funcao faz:
1 - pega um []byte coloca um nome nele
2 - Comprime ele de .gz para file e salva em um diretorio o arquivo (unico)
*/
func DecompressGzipToSaveOneFile(gzippedData []byte, outputPath string) error {
	// Criar o buffer de leitura
	byteReader := bytes.NewReader(gzippedData)

	// Criar o leitor GZIP
	gzipReader, err := gzip.NewReader(byteReader)
	if err != nil {
		return fmt.Errorf("erro ao criar leitor GZIP: %w", err)
	}
	defer gzipReader.Close()

	// Criar o arquivo de saída
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de saída: %w", err)
	}
	defer outFile.Close()

	// Copiar os dados descomprimidos para o arquivo
	_, err = io.Copy(outFile, gzipReader)
	if err != nil {
		return fmt.Errorf("erro ao gravar dados descomprimidos no arquivo: %w", err)
	}

	return nil
}

/*
DecompressGzipToBytes
Esta funcao faz:
1 - pega um []byte coloca um nome nele
2 - Comprime ele de .gz para file e devolve o byte dele assim voce pode optar por
salvar em disco ou ate salvar em banco
3 - retorna apenas o byte de um unico arquivo!
*/
func DecompressGzipToBytes(gzippedData []byte) ([]byte, error) {
	// Criar o buffer de leitura a partir dos dados comprimidos
	byteReader := bytes.NewReader(gzippedData)

	// Criar o leitor GZIP
	gzipReader, err := gzip.NewReader(byteReader)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar leitor GZIP: %w", err)
	}
	defer gzipReader.Close()

	// Buffer para armazenar os dados descomprimidos
	var decompressedData bytes.Buffer

	// Copiar os dados descomprimidos para o buffer
	_, err = io.Copy(&decompressedData, gzipReader)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler dados descomprimidos: %w", err)
	}

	// Retornar o conteúdo descomprimido como []byte
	return decompressedData.Bytes(), nil
}
