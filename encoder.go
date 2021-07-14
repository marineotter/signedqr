package signedqr

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func Encode(data string, secretKeyPath string) *gozxing.BitMatrix {
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
	signBase64 := base64.StdEncoding.EncodeToString(signature)
	fmt.Printf("hexstr: %d\n", len(signStr))
	fmt.Printf("base64: %d\n", len(signBase64))

	fmt.Printf("enc: hash: %x\n", hashed[:])
	fmt.Printf("enc: sign: %x\n", signBase64)

	signedData := data + signBase64

	// Generate QR Code
	w := qrcode.NewQRCodeWriter()
	hints := make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_MARGIN] = 0
	hints[gozxing.EncodeHintType_ERROR_CORRECTION] = "H"
	qr, err := w.Encode(signedData, gozxing.BarcodeFormat_QR_CODE, 512, 512, hints)

	return qr
}
