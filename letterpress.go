package gocarina

import (
	"image"
	_ "image/png" // register PNG format
	"log"
	"os"
)

// describes the geometry of the letterpress board source images
const (
	LetterpressTilesAcross    = 5
	LetterpressTilesDown      = 5
	LetterpressTilePixels     = 128
	LetterpressHeightOffset   = 496
	LetterPressExpectedWidth  = LetterpressTilesAcross * LetterpressTilePixels
	LetterpressExpectedHeight = 1136
)

// Board represents a Letterpress gameboard image
type Board struct {
	img   image.Image
	Tiles []*Tile
}

func ReadKnownBoard(file string, letters []rune) *Board {
	b := &Board{}
	b.img = ReadImage(file)
	images := b.crop()
	for i, img := range images {
		tile := NewTile(letters[i], img)
		b.Tiles = append(b.Tiles, tile)
	}

	return b
}

func ReadUnknownBoard(file string) *Board {
	b := &Board{}
	b.img = ReadImage(file)
	images := b.crop()
	for _, img := range images {
		tile := NewTile('?', img)
		b.Tiles = append(b.Tiles, tile)
	}

	return b
}

func ReadKnownBoards() map[rune]*Tile {
	result := make(map[rune]*Tile)

	letters := []rune{
		'P', 'R', 'B', 'R', 'Z',
		'T', 'A', 'V', 'Z', 'R',
		'B', 'D', 'A', 'K', 'Y',
		'G', 'I', 'G', 'K', 'F',
		'R', 'Y', 'S', 'J', 'V',
	}

	b := ReadKnownBoard("board-images/board1.png", letters)
	for _, tile := range b.Tiles {
		result[tile.Letter] = tile
	}

	letters = []rune{
		'Q', 'D', 'F', 'P', 'M',
		'N', 'E', 'E', 'S', 'I',
		'A', 'W', 'F', 'M', 'L',
		'F', 'R', 'P', 'T', 'T',
		'K', 'C', 'S', 'S', 'Y',
	}
	b = ReadKnownBoard("board-images/board2.png", letters)
	for _, tile := range b.Tiles {
		result[tile.Letter] = tile
	}

	letters = []rune{
		'L', 'H', 'F', 'L', 'M',
		'R', 'V', 'P', 'U', 'K',
		'V', 'O', 'E', 'E', 'X',
		'I', 'N', 'R', 'I', 'T',
		'V', 'N', 'S', 'I', 'Q',
	}
	b = ReadKnownBoard("board-images/board3.png", letters)
	for _, tile := range b.Tiles {
		result[tile.Letter] = tile
	}

	return result
}

func ReadImage(file string) image.Image {
	infile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

// crops a letterpress screen grab into a slice of tile images, one per letter.
//
func (b *Board) crop() (result []image.Image) {
	if b.img.Bounds().Dx() != LetterPressExpectedWidth || b.img.Bounds().Dy() != LetterpressExpectedHeight {
		log.Printf("Scaling...\n")
		b.img = Scale(b.img, image.Rect(0, 0, LetterPressExpectedWidth, LetterpressExpectedHeight))
	}

	yOffset := LetterpressHeightOffset
	border := 1

	for i := 0; i < LetterpressTilesDown; i++ {
		xOffset := 0

		for j := 0; j < LetterpressTilesAcross; j++ {
			tileRect := image.Rect(xOffset+border, yOffset+border, xOffset+LetterpressTilePixels-border, yOffset+LetterpressTilePixels-border)

			tile := b.img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(tileRect)

			result = append(result, tile)

			xOffset += LetterpressTilePixels
		}

		yOffset += LetterpressTilePixels
	}

	return
}
