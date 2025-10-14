package lib

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"log"

	"github.com/skip2/go-qrcode"
)

func GerarQrCodeJpegBase64(pixCode string) string {

	// Gerar o QR Code em alta resolução
	qrCode, err := qrcode.New(pixCode, qrcode.High)
	if err != nil {
		log.Fatalf("Erro ao gerar QR Code: %v", err)
	}

	// Criar um buffer para salvar a imagem como JPEG
	var buf bytes.Buffer

	// Codificar o QR Code para formato JPEG
	err = jpeg.Encode(&buf, qrCode.Image(256), nil) // 256 é o tamanho da imagem
	if err != nil {
		log.Fatalf("Erro ao codificar JPEG: %v", err)
	}

	// Converter a imagem em base64
	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Image
}
