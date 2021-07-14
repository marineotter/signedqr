package signedqr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestSuccess(t *testing.T) {
	ctx := time.Now().Format("2006-0102-150405") + "TestSuccess"

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

func TestDirtyCode(t *testing.T) {
	ctx := time.Now().Format("2006-0102-150405") + "TestDirtyCode"

	GenerateKeyPair(".", ctx)

	{
		qr := Encode("a", fmt.Sprintf("./%s.key", ctx))
		file, err1 := os.Create(fmt.Sprintf("%s.png", ctx))
		if err1 != nil {
			t.Fatal(err1)
		}
		buf := new(bytes.Buffer)
		_ = png.Encode(buf, qr)
		_, err2 := file.Write(buf.Bytes())
		file.Close()
		if err2 != nil {
			t.Fatal(err2)
		}
	}

	{
		// Stain
		file, err := os.Open(fmt.Sprintf("%s.png", ctx))
		img, _, err := image.Decode(file)
		outRect := image.Rectangle{image.Pt(0, 0), img.Bounds().Size()}
		out := image.NewRGBA(outRect)
		dstRect := image.Rectangle{image.Pt(0, 0), img.Bounds().Size()}
		draw.Draw(out, dstRect, img, image.Pt(0, 0), draw.Src)

		{
			x0 := 256
			y0 := 256
			r := 200
			c := color.RGBA{255, 100, 100, 255}
			x, y, dx, dy := r-1, 0, 1, 1
			err := dx - (r * 2)

			for x > y {
				thickness := 8
				for i := -thickness; i < thickness; i++ {
					for j := -thickness; j < thickness; j++ {
						out.Set(x0+x+i, y0+y+j, c)
						out.Set(x0+y+j, y0+x+i, c)
						out.Set(x0-y+j, y0+x+i, c)
						out.Set(x0-x+i, y0+y+j, c)
						out.Set(x0-x+i, y0-y+j, c)
						out.Set(x0-y+j, y0-x+i, c)
						out.Set(x0+y+j, y0-x+i, c)
						out.Set(x0+x+i, y0-y+j, c)
					}
				}

				if err <= 0 {
					y++
					err += dy
					dy += 2
				}
				if err > 0 {
					x--
					dx += 2
					err += dx - (r * 2)
				}
			}
		}
		file, err1 := os.Create(fmt.Sprintf("%s.png", ctx))
		if err1 != nil {
			t.Fatal(err1)
		}
		buf := new(bytes.Buffer)
		_ = png.Encode(buf, out)
		file.Write(buf.Bytes())
		file.Close()

		if err != nil {
			t.Fatal(err)
		}
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
