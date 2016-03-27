package gocarina

import (
	"image"
	_ "image/png" // register PNG format
	"log"
	"math"
	"os"
	"fmt"
	"image/png"
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


// describes the *target* geometry of the tiles, after we have sampled them down
const (
	TileTargetWidth   = 16
	TileTargetHeight  = 16
)

// populate map with reference images from letterpress game boards.

func ProcessGameBoards() map[rune]image.Image {
	result := make(map[rune]image.Image)

	img := ReadImage("board-images/board1.png")
	tiles := CropGameboard(img)
	tiles = DownsampleTiles(tiles)

	// TODO: delete this. it's for debugging the downsampled tiles
	for i, tile := range tiles {
		toFile, err := os.Create(fmt.Sprintf("downsampled%02d.png", i + 1))
		if err != nil {
			log.Fatal(err)
		}
		defer toFile.Close()

		err = png.Encode(toFile, tile)
		if err != nil {
			log.Fatal(err)
		}
	}


	runes := []rune{
		'P', 'R', 'B', 'R', 'Z',
		'T', 'A', 'V', 'Z', 'R',
		'B', 'D', 'A', 'K', 'Y',
		'G', 'I', 'G', 'K', 'F',
		'R', 'Y', 'S', 'J', 'V',
	}

	for i, r := range runes {
		result[r] = tiles[i]
	}

	img = ReadImage("board-images/board2.png")
	tiles = CropGameboard(img)
	tiles = DownsampleTiles(tiles)
	runes = []rune{
		'Q', 'D', 'F', 'P', 'M',
		'N', 'E', 'E', 'S', 'I',
		'A', 'W', 'F', 'M', 'L',
		'F', 'R', 'P', 'T', 'T',
		'K', 'C', 'S', 'S', 'Y',
	}

	for i, r := range runes {
		result[r] = tiles[i]
	}

	img = ReadImage("board-images/board3.png")
	tiles = CropGameboard(img)
	tiles = DownsampleTiles(tiles)
	runes = []rune{
		'L', 'H', 'F', 'L', 'M',
		'R', 'V', 'P', 'U', 'K',
		'V', 'O', 'E', 'E', 'X',
		'I', 'N', 'R', 'I', 'T',
		'V', 'N', 'S', 'I', 'Q',
	}

	for i, r := range runes {
		result[r] = tiles[i]
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

// CropGameboard crops a letterpress screen grab into a slice of tile images, one per letter.
//
// TODO: no bounding_box
func CropGameboard(img image.Image) (result []image.Image) {
	if img.Bounds().Dx() != LetterPressExpectedWidth || img.Bounds().Dy() != LetterpressExpectedHeight {
		log.Printf("Scaling...\n")
		img = Scale(img, image.Rect(0, 0, LetterPressExpectedWidth, LetterpressExpectedHeight))
	}

	img = BlackWhiteImage(img)

	yOffset := LetterpressHeightOffset
	border := 1

	for i := 0; i < LetterpressTilesDown; i++ {
		xOffset := 0

		for j := 0; j < LetterpressTilesAcross; j++ {
			tileRect := image.Rect(xOffset+border, yOffset+border, xOffset+LetterpressTilePixels-border, yOffset+LetterpressTilePixels-border)

			tile := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(tileRect)

			result = append(result, tile)

			xOffset += LetterpressTilePixels
		}

		yOffset += LetterpressTilePixels
	}

	return
}

// Scale scales the src image to the given rectangle using Nearest Neighbor
func Scale(srcImg image.Image, r image.Rectangle) image.Image {
	dstImg := image.NewRGBA(r)

	sw := srcImg.Bounds().Dx()
	sh := srcImg.Bounds().Dy()

	dw := dstImg.Bounds().Dx()
	dh := dstImg.Bounds().Dy()

	xAspect := float64(sw) / float64(dw)
	yAspect := float64(sh) / float64(dh)

	for y := 0; y < dh; y++ {
		for x := 0; x < dw; x++ {
			srcX := int(math.Floor(float64(x) * xAspect))
			srcY := int(math.Floor(float64(y) * yAspect))
			pix := srcImg.At(srcX, srcY)
			log.Printf("x: %d, srcX: %d, y: %d, srcY: %d, color: %+v", x, srcX, y, srcY, pix)
			dstImg.Set(x, y, pix)
		}
	}

	return dstImg
}


func DownsampleTiles(tiles []image.Image) (result []image.Image) {
	rect := image.Rect(0, 0, TileTargetWidth, TileTargetHeight)

	for _, tile := range tiles {
		downSampled := Scale(tile, rect)
		result = append(result, downSampled)
	}

	return
}
