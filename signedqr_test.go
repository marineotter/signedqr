package signedqr

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestExampleSuccess(t *testing.T) {
	ctx := time.Now().Format("2006-0102-150405")

	GenerateKeyPair(".", ctx)

	qr := Encode("Sample string", fmt.Sprintf("./%s.key", ctx))
	file, err1 := os.Create(fmt.Sprintf("%s.png", ctx))
	if err1 != nil {
		t.Fatal(err1)
	}
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, qr)
	_, err2 := file.Write(buf.Bytes())
	if err2 != nil {
		t.Fatal(err2)
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s.png", ctx))
	if err != nil {
		panic(err)
	}
	res, err := Decode(data, fmt.Sprintf("./%s.pub", ctx))
	println(res)
}
