package scene_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"github.com/d-s-hoffmann/go-ui/internal/scene"
)

// scenes returns a fresh instance of every scene. New scenes must be added
// here to be covered by the contract tests.
func scenes() []scene.Scene {
	return []scene.Scene{scene.Plasma{}, scene.Julia{}, &scene.Lissajous{}}
}

var sizes = []image.Point{{320, 200}, {390, 844}, {7, 3}, {1, 1}, {5, 0}, {0, 0}}

// sentinel marks pixels a scene has not written. Scenes write opaque
// pixels only, so alpha 0 never occurs in their output.
var sentinel = color.RGBA{1, 2, 3, 0}

func render(s scene.Scene, size image.Point, t time.Duration) *image.RGBA {
	img := image.NewRGBA(image.Rectangle{Max: size})
	s.Render(img, t)
	return img
}

func TestNamesAreUniqueAndNonEmpty(t *testing.T) {
	seen := map[string]bool{}
	for _, s := range scenes() {
		name := s.Name()
		if name == "" || seen[name] {
			t.Errorf("name %q is empty or duplicated", name)
		}
		seen[name] = true
	}
}

// TestRenderWritesExactlyBounds renders into a sub-image of a larger
// canvas, so a non-zero origin and a stride wider than the frame are
// exercised, and checks that every pixel inside was written and none
// outside was touched.
func TestRenderWritesExactlyBounds(t *testing.T) {
	for _, s := range scenes() {
		for _, size := range sizes {
			t.Run(fmt.Sprintf("%s/%v", s.Name(), size), func(t *testing.T) {
				canvas := image.NewRGBA(image.Rect(0, 0, size.X+6, size.Y+6))
				draw.Draw(canvas, canvas.Bounds(), image.NewUniform(sentinel), image.Point{}, draw.Src)
				inner := image.Rectangle{Max: size}.Add(image.Pt(3, 2))
				s.Render(canvas.SubImage(inner).(*image.RGBA), time.Second)

				b := canvas.Bounds()
				for y := b.Min.Y; y < b.Max.Y; y++ {
					for x := b.Min.X; x < b.Max.X; x++ {
						c := canvas.RGBAAt(x, y)
						inside := image.Pt(x, y).In(inner)
						if inside && c.A != 255 {
							t.Fatalf("pixel (%d,%d) inside the frame was not written: %v", x, y, c)
						}
						if !inside && c != sentinel {
							t.Fatalf("pixel (%d,%d) outside the frame was overwritten: %v", x, y, c)
						}
					}
				}
			})
		}
	}
}

// TestRenderIsDeterministic checks that a frame depends only on size and
// time, including after a scene has rendered other sizes, which catches
// state carried over in reused buffers.
func TestRenderIsDeterministic(t *testing.T) {
	size := image.Pt(64, 40)
	at := 1500 * time.Millisecond
	for i, s := range scenes() {
		t.Run(s.Name(), func(t *testing.T) {
			want := render(scenes()[i], size, at)
			render(s, image.Pt(100, 100), 3*time.Second)
			render(s, image.Pt(10, 10), 0)
			if got := render(s, size, at); !bytes.Equal(got.Pix, want.Pix) {
				t.Error("frame differs from a fresh scene's frame for the same size and time")
			}
		})
	}
}

func TestRenderAnimates(t *testing.T) {
	size := image.Pt(64, 40)
	for _, s := range scenes() {
		t.Run(s.Name(), func(t *testing.T) {
			a, b := render(s, size, 0), render(s, size, time.Second)
			if bytes.Equal(a.Pix, b.Pix) {
				t.Error("frames at 0s and 1s are identical")
			}
		})
	}
}

func TestLissajousDrawsStrokeOnBackground(t *testing.T) {
	img := render(&scene.Lissajous{}, image.Pt(320, 200), 0)
	background := color.RGBA{0x0b, 0x09, 0x20, 0xff}
	stroke := color.RGBA{0x8f, 0x6b, 0xff, 0xff}

	// The curve keeps a margin from the edges, so corners are background.
	if c := img.RGBAAt(0, 0); c != background {
		t.Errorf("corner = %v, want background %v", c, background)
	}
	var full int
	for y := range 200 {
		for x := range 320 {
			if img.RGBAAt(x, y) == stroke {
				full++
			}
		}
	}
	if full == 0 {
		t.Error("no pixel is fully covered by the stroke")
	}
}

func BenchmarkRender(b *testing.B) {
	for _, s := range scenes() {
		b.Run(s.Name(), func(b *testing.B) {
			img := image.NewRGBA(image.Rect(0, 0, 320, 200))
			for i := 0; b.Loop(); i++ {
				s.Render(img, time.Duration(i)*16*time.Millisecond)
			}
		})
	}
}
