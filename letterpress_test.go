package gocarina

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
)

// no assertions, but this exercises the entire board -> tile process, and it's also useful to get
// debugging images to written to debug_output/**
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

// Again, no assertions here, but handy way to create a board image with noise. This is a way to convince yourself
// that the network is doing more than a bit-per-bit image comparison. By running the "noised" board through
// the recognizer, we can see how it does on an image that has had some of its pixels disturbed.
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
