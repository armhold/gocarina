package gocarina

import (
	"image"
	"math"
	"log"
	"image/png"
	"fmt"
	"os"
)

const (
	MinBoundingBoxPercent = 0.25

        // describes the *target* geometry of the tiles, after we have sampled them down
	TileTargetWidth   = 12
	TileTargetHeight  = 12
)


// Tile represents a lettered square from a Letterpress gameboard.
type Tile struct {
	Letter  rune        // the letter this tile represents, if known
	img     image.Image // the original tile image, prior to any scaling/downsampling
	Reduced image.Image // the tile in black and white, bounding-boxed, and scaled down
}


func NewTile(letter rune, img image.Image) (result *Tile) {
	result = &Tile{Letter: letter, img: img}
	result.Reduce(2)

	return
}

// Reduce the tile by converting to monochrome, applying a bounding box, and scaling to match the given size.
// The resulting image will be stored in t.reducedImage
func (t *Tile) Reduce(border int) {
	targetRect := image.Rect(0, 0, TileTargetWidth, TileTargetHeight)

	src := BlackWhiteImage(t.img)

	// find the bounding box for the character
	bbox := boundingBox(src, 2)

	// Only apply the bounding box if it's above some % of the width/height of original tile.
	// This is to avoid pathological cases for skinny letters like "I", which
	// would otherwise result in completely black tiles when bounded.

	if bbox.Bounds().Dx() >= int(MinBoundingBoxPercent * float64(t.img.Bounds().Dx())) &&
	   bbox.Bounds().Dy() >= int(MinBoundingBoxPercent * float64(t.img.Bounds().Dy())) {
		src = src.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(bbox)
	} else {
		log.Printf("rune: %c: skipping boundingbox: orig width: %d, boundbox width: %d", t.Letter, t.img.Bounds().Dx(), bbox.Dx())
	}

	t.Reduced = scale(src, targetRect)
}


// BoundingBox returns the minimum rectangle containing all non-white pixels in the source image.
func boundingBox(src image.Image, border int) image.Rectangle {
	min := src.Bounds().Min
	max := src.Bounds().Max

	leftX := func() int {
		for x := min.X; x < max.X; x++ {
			for y := min.Y; y < max.Y; y++ {
				c := src.At(x, y)
				if IsBlack(c) {
					return x - border
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
					return x + border
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
					return y - border
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
					return y + border
				}
			}
		}

		// no non-white pixels found
		return max.Y
	}

	return image.Rect(leftX(), topY(), rightX(), bottomY())
}


// Scale scales the src image to the given rectangle using Nearest Neighbor
func scale(src image.Image, r image.Rectangle) image.Image {
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

func (t *Tile) SaveReducedTile() {
	toFile, err := os.Create(fmt.Sprintf("reduced_%c.png", t.Letter))
	if err != nil {
		log.Fatal(err)
	}
	defer toFile.Close()

	err = png.Encode(toFile, t.Reduced)
	if err != nil {
		log.Fatal(err)
	}
}
