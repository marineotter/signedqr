package signedqr

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateKeyPair(directory string, prefix string) error {
	reader := rand.Reader
	key, err := ecdsa.GenerateKey(elliptic.P224(), reader)
	if err != nil {
		error.Error(err)
	}

	// Store privatekey
	{
		var privateFileName = filepath.Join(directory, fmt.Sprintf("%s.key", prefix))
		of1, err := os.Create(privateFileName)
		if err != nil {
			error.Error(err)
		}
		defer of1.Close()
		data, _ := x509.MarshalECPrivateKey(key)
		err = pem.Encode(of1, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: data,
		})
		of1.Close()
	}

	// Store publickey
	{
		var publicFileName = filepath.Join(directory, fmt.Sprintf("%s.pub", prefix))
		of2, err := os.Create(publicFileName)
		if err != nil {
			error.Error(err)
		}
		defer of2.Close()
		data, err := x509.MarshalPKIXPublicKey(key.Public())
		err = pem.Encode(of2, &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: data,
		})
		of2.Close()
	}

	return nil
}
