package signedqr

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestExampleSuccess(t *testing.T) {
	ctx := time.Now().Format("2006-0102-150405")

	GenerateKeyPair(".", ctx)

	pngBinaryData := Encode("Sample string", fmt.Sprintf("./%s.key", ctx))
	file, err1 := os.Create(fmt.Sprintf("%s.png", ctx))
	if err1 != nil {
		t.Fatal(err1)
	}
	_, err2 := file.Write(pngBinaryData)
	if err2 != nil {
		t.Fatal(err2)
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s.png", ctx))
	if err != nil {
		panic(err)
	}
	res, err := Decode(data, fmt.Sprintf("./%s.pub", ctx))
	print(res)
}
