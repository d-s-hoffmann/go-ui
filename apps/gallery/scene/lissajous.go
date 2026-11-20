package scene

import (
	"image"
	"image/color"
	"math"
	"time"
)

const (
	lissajousSteps  = 1400 // line segments per curve
	lissajousMargin = 6    // distance of the curve's extremes from the edge, in pixels
	lissajousWidth  = 1.5  // stroke width, in pixels
)

var (
	lissajousBackground = color.RGBA{0x0b, 0x09, 0x20, 0xff}
	lissajousStroke     = color.RGBA{0x8f, 0x6b, 0xff, 0xff}
)

// Lissajous is a 3:4 Lissajous curve with a drifting phase.
//
// The zero value is ready to use. It keeps a scratch buffer between
// frames, so it must not be used concurrently.
type Lissajous struct {
	cov []float32 // stroke coverage per pixel, 0 to 1
}

// Name implements [Scene].
func (*Lissajous) Name() string { return "Lissajous" }

// Render implements [Scene].
func (l *Lissajous) Render(dst *image.RGBA, t time.Duration) {
	b := dst.Bounds()
	w, h := b.Dx(), b.Dy()
	if cap(l.cov) < w*h {
		l.cov = make([]float32, w*h)
	}
	cov := l.cov[:w*h]
	clear(cov)

	s := t.Seconds()
	cx, cy := float64(w)/2, float64(h)/2
	ax, ay := cx-lissajousMargin, cy-lissajousMargin
	point := func(i int) (x, y float64) {
		u := 2 * math.Pi * float64(i) / lissajousSteps
		return cx + ax*math.Sin(3*u+s*0.9), cy + ay*math.Sin(4*u)
	}
	x0, y0 := point(0)
	for i := 1; i <= lissajousSteps; i++ {
		x1, y1 := point(i)
		strokeSegment(cov, w, h, x0, y0, x1, y1, lissajousWidth/2)
		x0, y0 = x1, y1
	}

	rows(dst, func(y int, row []byte) {
		for x, c := range cov[y*w : (y+1)*w] {
			set(row[x*4:], mix(lissajousBackground, lissajousStroke, c))
		}
	})
}

// strokeSegment rasterises the segment (x0,y0)-(x1,y1) with half-width r
// into the w×h coverage buffer cov. Coverage ramps linearly over one pixel
// around the stroke edge. Each pixel keeps its highest coverage, so joints
// and self-intersections don't get brighter than the stroke itself.
func strokeSegment(cov []float32, w, h int, x0, y0, x1, y1, r float64) {
	pad := r + 0.5
	minX := max(int(math.Floor(min(x0, x1)-pad)), 0)
	maxX := min(int(math.Ceil(max(x0, x1)+pad)), w-1)
	minY := max(int(math.Floor(min(y0, y1)-pad)), 0)
	maxY := min(int(math.Ceil(max(y0, y1)+pad)), h-1)

	dx, dy := x1-x0, y1-y0
	len2 := dx*dx + dy*dy
	for py := minY; py <= maxY; py++ {
		for px := minX; px <= maxX; px++ {
			// Distance from the pixel centre to the closest point on the segment.
			qx, qy := float64(px)+0.5-x0, float64(py)+0.5-y0
			k := 0.0
			if len2 > 0 {
				k = min(max((qx*dx+qy*dy)/len2, 0), 1)
			}
			d := math.Hypot(qx-k*dx, qy-k*dy)
			c := float32(min(max(pad-d, 0), 1))
			i := py*w + px
			cov[i] = max(cov[i], c)
		}
	}
}

// mix blends from a to b by k in [0, 1], returning an opaque colour.
func mix(a, b color.RGBA, k float32) color.RGBA {
	if k <= 0 {
		return a
	}
	ch := func(x, y uint8) uint8 {
		return uint8(float32(x) + (float32(y)-float32(x))*k + 0.5)
	}
	return color.RGBA{ch(a.R, b.R), ch(a.G, b.G), ch(a.B, b.B), 255}
}
