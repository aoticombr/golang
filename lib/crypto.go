package lib

import (
	"crypto"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"os"
)

/*
EncodeByte:
Recebe uma []byte como entrada e retorna a string codificada em base64
Imita funcao da unit Soap.EncdDecd do Delphi
*/
func EncodeByte(input []byte) string {
	encoded := base64.StdEncoding.EncodeToString(input)
	return encoded
}

/*
EncodeString:
Recebe uma string como entrada e retorna a string codificada em base64
Imita funcao da unit Soap.EncdDecd do Delphi
*/
func EncodeString(input string) string {
	encoded := EncodeByte([]byte(input))
	return encoded
}

func Md5(src string) string {
	hashedPassword := md5.Sum([]byte(src))
	passwordMD5 := hex.EncodeToString(hashedPassword[:])
	return passwordMD5
}

// Função para carregar a chave pública a partir de um arquivo PEM
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("falha ao carregar a chave pública")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubKey.(*rsa.PublicKey), nil
}

// Função para carregar a chave privada a partir de um arquivo PEM
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privKeyPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("falha ao carregar a chave privada")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

// Função para criptografar dados com a chave pública RSA
func EncryptData(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedData, err := rsa.EncryptOAEP(
		sha256.New(),
		rand2.Reader,
		publicKey,
		data,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

// Função para descriptografar dados com a chave privada RSA
func DecryptData(privateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	decryptedData, err := rsa.DecryptOAEP(
		sha256.New(),
		rand2.Reader,
		privateKey,
		encryptedData,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

// Função para assinar dados com a chave privada do remetente
func SignData(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand2.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}
