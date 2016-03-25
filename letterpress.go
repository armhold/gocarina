package gocarina

import (
	"image"
	"log"
	"math"
)

const (
	LetterpressTilesAcross    = 5
	LetterpressTilesDown      = 5
	LetterpressTilePixels     = 128
	LetterpressHeightOffset   = 496
	LetterPressExpectedWidth  = LetterpressTilesAcross * LetterpressTilePixels
	LetterpressExpectedHeight = 1136
)

// TODO: no bounding_box
func Crop(img image.Image) (result [][]image.Image) {
	if img.Bounds().Dx() != LetterPressExpectedWidth || img.Bounds().Dy() != LetterpressExpectedHeight {
		log.Printf("Scaling...\n")
		img = Scale(img, image.Rect(0, 0, LetterPressExpectedWidth, LetterpressExpectedHeight))
	}

	img = BlackWhiteImage(img)

	yOffset := LetterpressHeightOffset
	border := 1

	for i := 0; i < LetterpressTilesDown; i++ {
		xOffset := 0
		var row []image.Image

		for j := 0; j < LetterpressTilesAcross; j++ {
			tileRect := image.Rect(xOffset+border, yOffset+border, xOffset+LetterpressTilePixels-border, yOffset+LetterpressTilePixels-border)

			tile := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(tileRect)

			row = append(row, tile)

			xOffset += LetterpressTilePixels
		}
		result = append(result, row)

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
			dstImg.Set(x, y, pix)
		}
	}

	return dstImg
}
