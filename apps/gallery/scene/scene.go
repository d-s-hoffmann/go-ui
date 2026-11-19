// Package scene implements animated scenes that render into an
// [image.RGBA]. It must not depend on any browser API, so scenes can be
// rendered and tested natively.
package scene

import (
	"image"
	"time"
)

// A Scene draws one animation frame at a time.
//
// Render must write every pixel of dst.Bounds() and nothing outside it,
// and must produce the same pixels for the same size and time. The size of
// dst may change between calls.
type Scene interface {
	// Name is a short human-readable label, e.g. for a menu entry.
	Name() string
	// Render draws the frame at time t, measured from the animation start.
	Render(dst *image.RGBA, t time.Duration)
}

// rows calls fn once per row of dst with that row's pixel bytes, four per
// pixel in RGBA order.
func rows(dst *image.RGBA, fn func(y int, row []byte)) {
	b := dst.Bounds()
	n := b.Dx() * 4
	for y := range b.Dy() {
		off := dst.PixOffset(b.Min.X, b.Min.Y+y)
		fn(y, dst.Pix[off:off+n:off+n])
	}
}
