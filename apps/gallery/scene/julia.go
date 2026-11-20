package scene

import (
	"image"
	"image/color"
	"math"
	"time"
)

// juliaIter is the escape-time iteration limit.
const juliaIter = 60

// juliaInside is the colour of points that never escape.
var juliaInside = color.RGBA{12, 9, 36, 255}

// Julia is a Julia set whose constant slowly orbits the origin.
type Julia struct{}

// Name implements [Scene].
func (Julia) Name() string { return "Julia set" }

// Render implements [Scene].
func (Julia) Render(dst *image.RGBA, t time.Duration) {
	angle := t.Seconds() * 0.15
	cr := math.Cos(angle) * 0.7885
	ci := math.Sin(angle) * 0.7885

	b := dst.Bounds()
	w, h := float64(b.Dx()), float64(b.Dy())
	halfW, halfH := float64(b.Dx()/2), float64(b.Dy()/2)
	rows(dst, func(y int, row []byte) {
		zy0 := (float64(y) - halfH) / (h / 2.6)
		for x := range len(row) / 4 {
			zx := (float64(x) - halfW) / (w / 3.2)
			zy := zy0
			n := 0
			var zx2, zy2 float64
			for ; n < juliaIter; n++ {
				zx2, zy2 = zx*zx, zy*zy
				if zx2+zy2 > 4 {
					break
				}
				zy = 2*zx*zy + ci
				zx = zx2 - zy2 + cr
			}
			set(row[x*4:], juliaColor(n, zx2+zy2))
		}
	})
}

// juliaColor maps an escape count and the final |z|² to a colour, using
// the smooth (normalised) iteration count to avoid banding.
func juliaColor(n int, abs2 float64) color.RGBA {
	if n >= juliaIter {
		return juliaInside
	}
	mu := max(float64(n)+1-math.Log2(math.Log(math.Sqrt(abs2))), 0)
	return rainbow[int(mu*16)%len(rainbow)]
}
