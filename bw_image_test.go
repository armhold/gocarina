package gocarina

import (
	"image"
	"testing"
	"os"
	"image/color"
	_ "image/png"  // register PNG format
)

func TestBlackWhiteImage(t *testing.T) {
	infile, err := os.Open("ocarina.png")
	if err != nil {
		t.Fatal(err)
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		t.Fatal(err)
	}

	bwImg := BlackWhiteImage(img)

	for x := 0; x < bwImg.Bounds().Dx(); x++ {
		for y := 0; y < bwImg.Bounds().Dy(); y++ {
			c := bwImg.At(x, y)
			if ! isBlackOrWhite(c) {
				t.Fatal("not black or white: %+v", c)
			}
		}
	}

	//toFile, err := os.Create("ocarina_bw.png")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer toFile.Close()
	//
	//err = png.Encode(toFile, bwImg)
	//if err != nil {
	//	t.Fatal(err)
	//}
}

func isBlackOrWhite(c color.Color) bool {
	br, bg, bb, ba := color.Black.RGBA()
	wr, wg, wb, wa := color.White.RGBA()

	r, g, b, a := c.RGBA()

	if r == br && g == bg && b == bb && a == ba {
		return true
	}

	return r == wr && g == wg && b == wb && a == wa
}
