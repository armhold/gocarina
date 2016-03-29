package gocarina

import (
	"image"
	"image/color"
	"testing"
)

// just document behavior of image.Rectangle
func TestGeometry(t *testing.T) {
	r := image.Rect(0, 0, 16, 16)

	if r.Dx() != 16 {
		t.Fatalf("expected width to be 16")
	}

	if r.Dy() != 16 {
		t.Fatalf("expected height to be 16")
	}

	if r.Min.X != 0 {
		t.Fatalf("expected starting X coord to be 0")
	}

	// the max is not intended to be included in the range
	if r.Max.X != 16 {
		t.Fatalf("expected ending X coord to be 16")
	}

	if r.Min.Y != 0 {
		t.Fatalf("expected starting Y coord to be 0")
	}

	// the max is not intended to be included in the range
	if r.Max.Y != 16 {
		t.Fatalf("expected ending Y coord to be 16")
	}
}

func TestBoundingBox(t *testing.T) {
	r := image.Rect(0, 0, 16, 16)
	img := image.NewRGBA(r)

	// top left
	img.Set(3, 3, color.Black)

	// bottom right
	img.Set(12, 12, color.Black)

	bbox := BoundingBox(img, 0)
	assertWidth(bbox, 10, t)
	assertHeight(bbox, 10, t)

	// now test with border
	bbox = BoundingBox(img, 1)
	assertWidth(bbox, 12, t)
	assertHeight(bbox, 12, t)
}

func assertWidth(rect image.Rectangle, w int, t *testing.T) {
	if rect.Bounds().Dx() != w {
		t.Fatalf("expected rect.Bounds().Dx() to be %d, was: %d", w, rect.Bounds().Dx())
	}
}

func assertHeight(rect image.Rectangle, h int, t *testing.T) {
	if rect.Bounds().Dy() != h {
		t.Fatalf("expected rect.Bounds().Dy() to be %d, was: %d", h, rect.Bounds().Dy())
	}
}
