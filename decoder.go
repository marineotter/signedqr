package signedqr

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
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
	genericPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	key := genericPublicKey.(*ecdsa.PublicKey)
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
	keyStrSize := 0
	for i := 0; i < len(resultStr); i++ {
		if int(resultStr[len(resultStr)-1-i]) == 0 {
			keyStrSize = i
			break
		}
	}
	sign := result.String()[len(resultStr)-keyStrSize:]
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	data := result.String()[:len(resultStr)-keyStrSize-1]
	message := []byte(data)
	hashed := sha256.Sum224(message)
	fmt.Printf("dec: hash: %x\n", hashed[:])
	fmt.Printf("dec: sign: %x\n", signBytes)

	isOk := ecdsa.VerifyASN1(key, hashed[:], signBytes)
	if !isOk {
		return "", errors.New("Verification failed.")
	}
	return data, nil
}
