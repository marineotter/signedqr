package main

import (
	"fmt"
	"os"

	"github.com/marineotter/signedqr"
)

func sign(ctx string) {
	pngBinaryData := signedqr.Encode("てすとでーたです", fmt.Sprintf("./%s.key", ctx))

	file, err1 := os.Create(fmt.Sprintf("%s.png", ctx))
	if err1 != nil {
		fmt.Println("file create err:", err1)
		return
	}

	// バイナリデータをファイルに書き込み
	_, err2 := file.Write(pngBinaryData)
	if err2 != nil {
		fmt.Println("file write err:", err2)
		return
	}
}
