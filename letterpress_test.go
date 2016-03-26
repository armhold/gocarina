package gocarina

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
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

	tiles := CropGameboard(lp)


	for i, tile := range tiles {
		toFile, err := os.Create(fmt.Sprintf("tile%02d.png", i + 1))
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
