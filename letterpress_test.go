package gocarina

import (
	"fmt"
	"image/png"
	"os"
	"testing"
)

func TestCrop(t *testing.T) {
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
