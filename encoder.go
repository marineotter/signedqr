package signedqr

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"image/png"
	"io/ioutil"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func Encode(data string, secretKeyPath string) []byte {
	// Generate SignedString
	keydata, err := ioutil.ReadFile(secretKeyPath)
	if err != nil {
		print(err.Error())
		return nil
	}
	block, _ := pem.Decode(keydata)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		print(err.Error())
		return nil
	}
	message := []byte(data)
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	signStr := fmt.Sprintf("%x", signature)[:8]
	signedData := signStr + data

	// Generate QR Code
	w := qrcode.NewQRCodeWriter()
	qr, err := w.EncodeWithoutHint(signedData, gozxing.BarcodeFormat_QR_CODE, 512, 512)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, qr)

	return buf.Bytes()
}
