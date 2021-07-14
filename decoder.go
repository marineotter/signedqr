package signedqr

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"image"
	"io/ioutil"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func Decode(imgData []byte, publicKeyPath string) (string, error) {
	// Generate SignedString
	keydata, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		print(err.Error())
		return "", err
	}
	block, _ := pem.Decode(keydata)
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		print(err.Error())
		return "", err
	}

	// Read from QR Code
	img, _, _ := image.Decode(bytes.NewReader(imgData))
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if result == nil {
		return "", err
	}

	// Verify
	resultStr := result.String()
	sign := result.String()[len(resultStr)-88:]
	signBytes := make([]byte, base64.StdEncoding.DecodedLen(len(sign)))
	n, err := base64.StdEncoding.Decode(signBytes, []byte(sign))
	signBytes = signBytes[:n]
	data := result.String()[:len(resultStr)-88]
	print(data)
	message := []byte(data)
	hashed := sha256.Sum256(message)
	fmt.Printf("dec: hash: %x\n", hashed[:])
	fmt.Printf("dec: sign: %x\n", signBytes)

	err2 := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signBytes)
	if err2 != nil {
		print(err2.Error())
		return "", err2
	}
	return data, nil
}
