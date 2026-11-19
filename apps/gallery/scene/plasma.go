package scene

import (
	"image"
	"math"
	"time"
)

// Plasma is an animated plasma made of overlapping sine waves.
type Plasma struct{}

// Name implements [Scene].
func (Plasma) Name() string { return "Plasma" }

// Render implements [Scene].
func (Plasma) Render(dst *image.RGBA, t time.Duration) {
	s := t.Seconds()
	b := dst.Bounds()
	cx, cy := float64(b.Dx()/2), float64(b.Dy()/2)
	rows(dst, func(y int, row []byte) {
		fy := float64(y)
		for x := range len(row) / 4 {
			fx := float64(x)
			v := math.Sin(fx/18+s) +
				math.Sin(fy/9-s*1.4) +
				math.Sin((fx+fy)/26+s*0.8) +
				math.Sin(math.Hypot(fx-cx, fy-cy)/14-s)
			idx := min(max(int((v+4)*32), 0), len(rainbow)-1)
			set(row[x*4:], rainbow[idx])
		}
	})
}
