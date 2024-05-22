package goheif

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/nfnt/resize"
)

func TestFormatRegistered(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/camel.heic")
	if err != nil {
		t.Fatal(err)
	}

	img, dec, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("unable to decode heic image: %s", err)
	}

	if got, want := dec, "heic"; got != want {
		t.Errorf("unexpected decoder: got %s, want %s", got, want)
	}

	if w, h := img.Bounds().Dx(), img.Bounds().Dy(); w != 1596 || h != 1064 {
		t.Errorf("unexpected decoded image size: got %dx%d, want 1596x1064", w, h)
	}
}

func BenchmarkSafeEncoding(b *testing.B) {
	benchEncoding(b, true)
}

func BenchmarkRegularEncoding(b *testing.B) {
	benchEncoding(b, false)
}

func benchEncoding(b *testing.B, safe bool) {
	b.Helper()

	currentSetting := SafeEncoding
	defer func() {
		SafeEncoding = currentSetting
	}()
	SafeEncoding = safe

	f, err := ioutil.ReadFile("testdata/camel.heic")
	if err != nil {
		b.Fatal(err)
	}
	r := bytes.NewReader(f)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Decode(r)
		r.Seek(0, io.SeekStart)
	}
}

func TestHEIC(t *testing.T) {
	fileName := "./testdata/camel.heic" //IMG_20240522_142801.HEIC
	toFile := "mi.png"
	file, _ := os.Open(fileName)
	file2, _ := os.Open(toFile)
	defer file.Close()
	defer file2.Close()

	img0, err := Decode(file)
	if err != nil {
		panic(err)
	}

	newImg := resize.Resize(100, 100, img0, resize.Bilinear)

	t.Log(newImg.Bounds().String())

	ni := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			ni.Set(i, j, newImg.At(i, j))
		}
	}

	err = png.Encode(file2, ni)
	if err != nil {
		panic(err)
	}
	t.Log("success")
}
