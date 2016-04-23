package gocarina

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
)

// not a true test, but writes debugging images to debug_output/**
func TestReadKnownBoards(t *testing.T) {
	m := ReadKnownBoards()

	for letter, tile := range m {
		toFile, err := os.Create(fmt.Sprintf("debug_output/tile_%c.png", letter))
		if err != nil {
			t.Fatal(err)
		}
		defer toFile.Close()

		err = png.Encode(toFile, tile.img)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestNoise(t *testing.T) {
	infile, err := os.Open("board-images/board1.png")
	if err != nil {
		t.Fatal(err)
	}
	defer infile.Close()

	srcImg, _, err := image.Decode(infile)
	if err != nil {
		t.Fatal(err)
	}

	noiseyImg := ConvertToRGBA(srcImg)
	AddNoise(noiseyImg)

	toFile, err := os.Create("debug_output/board1-noise.png")
	if err != nil {
		t.Fatal(err)
	}
	defer toFile.Close()

	err = png.Encode(toFile, noiseyImg)
	if err != nil {
		t.Fatal(err)
	}

}
