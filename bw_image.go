package gocarina

import (
	"image"
	"image/color"
)

var (
	bwPalette      = color.Palette([]color.Color{color.Black, color.White})
	br, bg, bb, ba = color.Black.RGBA()
	wr, wg, wb, wa = color.White.RGBA()
)

// Converted uses a Black & White color model to quantize images to black & white.
// credit to Hjulle: http://stackoverflow.com/a/17076395/93995
//
type Converted struct {
	Img image.Image
	Mod color.Model
}

func (c *Converted) ColorModel() color.Model {
	return c.Mod
}

func (c *Converted) Bounds() image.Rectangle {
	return c.Img.Bounds()
}

// At forwards the call to the original image, then quantizes to Black or White by
// applying a threshold.
func (c *Converted) At(x, y int) color.Color {
	r, g, b, _ := c.Img.At(x, y).RGBA()

	combined := r + g + b

	if combined < 50000 {
		return color.Black
	}

	return color.White
}

func (c *Converted) SubImage(r image.Rectangle) image.Image {
	sub := c.Img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(r)

	// preserve the B&W color model
	return &Converted{sub, bwPalette}
}

func BlackWhiteImage(img image.Image) image.Image {
	return &Converted{img, bwPalette}
}

func IsBlack(c color.Color) bool {
	r, g, b, a := c.RGBA()

	return r == br && g == bg && b == bb && a == ba

	return r == wr && g == wg && b == wb && a == wa
}

func IsWhite(c color.Color) bool {
	r, g, b, a := c.RGBA()

	return r == wr && g == wg && b == wb && a == wa
}
