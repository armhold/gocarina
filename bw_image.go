package gocarina

import (
	"image"
	"image/color"
)

var (
	bwPalette = color.Palette([]color.Color{color.Black, color.White})
)


// Converted uses a Black & White color model to quantize images to black & white.
// credit to Hjulle: http://stackoverflow.com/a/17076395/93995
//
type Converted struct {
	Img image.Image
	Mod color.Model
}

func (c *Converted) ColorModel() color.Model{
	return c.Mod
}

func (c *Converted) Bounds() image.Rectangle{
	return c.Img.Bounds()
}

// At forwards the call to the original image and then asks the color model to convert it.
func (c *Converted) At(x, y int) color.Color{
	return c.Mod.Convert(c.Img.At(x,y))
}

func BlackWhiteImage(img image.Image) (image.Image) {
	return &Converted{img, bwPalette}
}

