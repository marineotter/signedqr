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

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
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
	signStr := fmt.Sprintf("%x", signature)

	fmt.Printf("enc: hash: %x\n", hashed[:])
	fmt.Printf("enc: sign: %x\n", signature)

	signedData := data + signStr

	// Generate QR Code

	qrData, _ := qr.Encode(signedData, qr.M, qr.Auto)
	qrCode, _ := barcode.Scale(qrData, 512, 512)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, qrCode)

	return buf.Bytes()
}
