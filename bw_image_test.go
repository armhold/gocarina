package gocarina

import (
	"image"
	"os"
	"testing"
)

func TestBlackWhiteImage(t *testing.T) {
	infile, err := os.Open("board-images/board1.png")
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
			if !(IsBlack(c) || IsWhite(c)) {
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
