package pgp

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/openpgp"
)

func DecodePGPKeyPass(encString, secretKeyring, passphrase string) string {

	//log.Println("Secret Keyring:", secretKeyring)
	//log.Println("Passphrase:", passphrase)

	// init some vars
	var entity *openpgp.Entity
	var entityList openpgp.EntityList

	//secretKeyring ja contem a chave privada
	keyringFileBuffer := bytes.NewBufferString(secretKeyring)

	entityList, err := openpgp.ReadArmoredKeyRing(keyringFileBuffer)
	//entityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return "erro"
	}
	entity = entityList[0]

	// Get the passphrase and read the private key.
	// Have not touched the encrypted string yet
	passphraseByte := []byte(passphrase)
	//log.Println("Decrypting private key using passphrase")
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}
	//log.Println("Finished decrypting private key using passphrase")

	// Decode the base64 string
	dec, err := base64.StdEncoding.DecodeString(encString)
	if err != nil {
		return "erro"
	}

	// Decrypt it with the contents of the private key
	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entityList, nil, nil)
	if err != nil {
		return "erro"
	}
	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "erro"
	}
	decStr := string(bytes)

	return decStr
}

func EncodePGPKey(secretString, secretKeyring string) (string, error) {

	keyringFileBuffer := bytes.NewBufferString(secretKeyring)
	entityList, err := openpgp.ReadArmoredKeyRing(keyringFileBuffer)
	if err != nil {
		return "", err
	}

	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte(secretString))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	// Encode to base64
	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}
	encStr := base64.StdEncoding.EncodeToString(bytes)

	return encStr, nil
}

func EncodePGPFile(secretString, NameFile string) (string, error) {

	keyringFileBuffer, err := os.Open(NameFile)
	if err != nil {
		return "", err
	}
	defer keyringFileBuffer.Close()

	keyringBytes, err := ioutil.ReadAll(keyringFileBuffer)
	if err != nil {
		return "", err
	}

	return EncodePGPKey(secretString, string(keyringBytes))
}

/*func EncodePGPFile(secretString, NameFile string) (string, error) {

	keyringFileBuffer, err := os.Open(NameFile)
	defer keyringFileBuffer.Close()

	entityList, err := openpgp.ReadArmoredKeyRing(keyringFileBuffer)
	if err != nil {
		return "", err
	}

	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte(secretString))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	// Encode to base64
	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}
	encStr := base64.StdEncoding.EncodeToString(bytes)

	return encStr, nil
}*/
