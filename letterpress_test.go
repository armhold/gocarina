package gocarina

import (
	"testing"
	"os"
	"image"
	"fmt"
	"image/png"
)


func TestCrop(t *testing.T) {
	infile, err := os.Open("ocarina.png")
	if err != nil {
		t.Fatal(err)
	}
	defer infile.Close()

	lp, _, err := image.Decode(infile)
	if err != nil {
		t.Fatal(err)
	}

	tiles := Crop(lp)

	tileCount := 1

	for _, row := range tiles {
		for _, tile := range row {
			toFile, err := os.Create(fmt.Sprintf("tile%02d.png", tileCount))
			if err != nil {
				t.Fatal(err)
			}
			defer toFile.Close()

			err = png.Encode(toFile, tile)
			if err != nil {
				t.Fatal(err)
			}

			tileCount++
		}
	}

}
