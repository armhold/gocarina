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
func Scale(src image.Image, r image.Rectangle) image.Image {
	dst := image.NewRGBA(r)

	sb := src.Bounds()
	db := dst.Bounds()

	for y := db.Min.Y; y < db.Max.Y; y++ {
		percentDownDest := float64(y) / float64(db.Dy())

		for x := db.Min.X; x < db.Max.X; x++ {

			percentAcrossDest := float64(x) / float64(db.Dx())

			srcX := int(math.Floor(percentAcrossDest * float64(sb.Dx())))
			srcY := int(math.Floor(percentDownDest * float64(sb.Dy())))

			pix := src.At(sb.Min.X + srcX, sb.Min.Y + srcY)
			dst.Set(x, y, pix)
		}
	}

	return dst
}


// BoundingBox returns the minimum rectangle containing all non-white pixels in the source image.
func BoundingBox(src image.Image) image.Rectangle {
	min := src.Bounds().Min
	max := src.Bounds().Max

	leftX := func() int {
		for x := min.X; x < max.X; x++ {
			for y := min.Y; y < max.Y; y++ {
				c := src.At(x, y)
				if IsBlack(c) {
					return x
				}
			}
		}

		// no non-white pixels found
		return min.X
	}

	rightX := func() int {
		for x := max.X - 1; x >= min.X; x-- {
			for y := min.Y; y < max.Y; y++ {
				c := src.At(x, y)
				if IsBlack(c) {
					return x
				}
			}
		}

		// no non-white pixels found
		return max.X
	}

	topY := func() int {
		for y := min.Y; y < max.Y; y++ {
			for x := min.X; x < max.X; x++ {
				c := src.At(x, y)
				if IsBlack(c) {
					return y
				}
			}
		}

		// no non-white pixels found
		return max.Y
	}

	bottomY := func() int {
		for y := max.Y - 1; y >= min.Y; y-- {
			for x := min.X; x < max.X; x++ {
				c := src.At(x, y)
				if IsBlack(c) {
					return y
				}
			}
		}

		// no non-white pixels found
		return max.Y
	}

	return image.Rect(leftX(), topY(), rightX(), bottomY())
}

func DownsampleTiles(tiles []image.Image) (result []image.Image) {
	rect := image.Rect(0, 0, TileTargetWidth, TileTargetHeight)

	for _, tile := range tiles {
		boundedRect := BoundingBox(tile)

		boundedImg := tile.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(boundedRect)

		downSampled := Scale(boundedImg, rect)
		result = append(result, downSampled)
	}

	return
}
