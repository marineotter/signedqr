package signedqr

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateKeyPair(directory string, prefix string) error {
	reader := rand.Reader
	bitSize := 512

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		error.Error(err)
	}
	// privatekeyの出力
	var privateFileName = filepath.Join(directory, fmt.Sprintf("%s.key", prefix))
	of1, err := os.Create(privateFileName)
	if err != nil {
		error.Error(err)
	}
	defer of1.Close()
	err = pem.Encode(of1, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	// privatekeyの出力
	var publicFileName = filepath.Join(directory, fmt.Sprintf("%s.pub", prefix))
	of2, err := os.Create(publicFileName)
	if err != nil {
		error.Error(err)
	}
	defer of2.Close()
	asn1Bytes, err := asn1.Marshal(key.PublicKey)
	err = pem.Encode(of2, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	})

	return nil
}
