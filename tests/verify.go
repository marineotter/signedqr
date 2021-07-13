package main

import (
	"fmt"
	"io/ioutil"

	"github.com/marineotter/signedqr"
)

func verify(ctx string) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s.png", ctx))
	if err != nil {
		panic(err)
	}
	signedqr.Decode(data, fmt.Sprintf("./%s.pub", ctx))
}
