package gocarina

import (
	"testing"
	"os"
	"image"
	"fmt"
	"image/png"
)


func TestCrop(t *testing.T) {
	infile, err := os.Open("board-images/board1.png")
	if err != nil {
		t.Fatal(err)
	}
	defer infile.Close()

	lp, _, err := image.Decode(infile)
	if err != nil {
		t.Fatal(err)
	}

	tiles := Crop(lp)

	n := 0

	for _, row := range tiles {
		for _, tile := range row {
			n++

			toFile, err := os.Create(fmt.Sprintf("tile%02d.png", n))
			if err != nil {
				t.Fatal(err)
			}
			defer toFile.Close()

			err = png.Encode(toFile, tile)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

}
