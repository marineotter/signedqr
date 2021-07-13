package signedqr

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
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
	result, _ := qrReader.Decode(bmp, nil)

	// Verify
	resultStr := result.String()
	sign := result.String()[len(resultStr)-512:]
	signBytes := make([]byte, hex.DecodedLen(len(sign)))
	n, err := hex.Decode(signBytes, []byte(sign))
	signBytes = signBytes[:n]
	data := result.String()[:len(resultStr)-512]
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
