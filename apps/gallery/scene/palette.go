package scene

import (
	"image/color"
	"math"
)

// rainbow is the shared colour palette of the pixel scenes.
var rainbow = newRainbow(256)

// newRainbow builds a bright HSV rainbow with n opaque entries.
func newRainbow(n int) []color.RGBA {
	pal := make([]color.RGBA, n)
	for i := range n {
		h := float64(i) / float64(n) * 6
		x := uint8(255 * (1 - math.Abs(math.Mod(h, 2)-1)))
		var r, g, b uint8
		switch int(h) % 6 {
		case 0:
			r, g, b = 255, x, 0
		case 1:
			r, g, b = x, 255, 0
		case 2:
			r, g, b = 0, 255, x
		case 3:
			r, g, b = 0, x, 255
		case 4:
			r, g, b = x, 0, 255
		default:
			r, g, b = 255, 0, x
		}
		pal[i] = color.RGBA{r, g, b, 255}
	}
	return pal
}

// set writes c into the four bytes of an RGBA pixel.
func set(px []byte, c color.RGBA) {
	px[0], px[1], px[2], px[3] = c.R, c.G, c.B, c.A
}
