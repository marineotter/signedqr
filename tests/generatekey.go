package main

import "github.com/marineotter/signedqr"

func generatekey(ctx string) {
	signedqr.GenerateKeyPair(".", ctx)
}
