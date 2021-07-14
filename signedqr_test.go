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

func TestSuccess(t *testing.T) {
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
		t.Fatal(err)
	}
	res, err := Decode(data, fmt.Sprintf("./%s.pub", ctx))
	if err != nil {
		t.Fatal(err)
	}
	println(res)
}

func TestVerifyWithAnotherKey(t *testing.T) {
	ctx_a := time.Now().Format("2006-0102-150405") + "-A"
	ctx_b := time.Now().Format("2006-0102-150405") + "-B"

	GenerateKeyPair(".", ctx_a)
	GenerateKeyPair(".", ctx_b)

	qr := Encode("Sample string", fmt.Sprintf("./%s.key", ctx_a))
	file, err1 := os.Create(fmt.Sprintf("%s.png", ctx_a))
	if err1 != nil {
		t.Fatal(err1)
	}
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, qr)
	_, err2 := file.Write(buf.Bytes())
	if err2 != nil {
		t.Fatal(err2)
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s.png", ctx_a))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Decode(data, fmt.Sprintf("./%s.pub", ctx_b))

	// It would be strange if I didn't get an error.
	if err == nil {
		t.Fatal(err)
	}
	println(res)
}
