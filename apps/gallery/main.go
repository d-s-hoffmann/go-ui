//go:build js && wasm

package main

import (
	"math"
	"time"

	"syscall/js"
)

const (
	width   = 320
	height  = 200
	maxIter = 60
)

var mode = "plasma"

func main() {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "stage")
	if canvas.IsUndefined() {
		println("no #stage canvas found")
		return
	}
	ctx := canvas.Call("getContext", "2d")

	img := ctx.Call("createImageData", width, height)
	pixels := make([]byte, width*height*4)
	buf := js.Global().Get("Uint8Array").New(len(pixels))
	pal := makePalette(256)

	type scene struct{ key, label string }
	scenes := []scene{
		{"plasma", "Plasma"},
		{"julia", "Julia set"},
		{"lissajous", "Lissajous"},
	}
	bar := doc.Call("getElementById", "controls")
	for _, sc := range scenes {
		key := sc.key
		btn := doc.Call("createElement", "button")
		btn.Set("textContent", sc.label)
		btn.Set("className", "btn")
		onClick := js.FuncOf(func(this js.Value, _ []js.Value) interface{} {
			mode = key
			return nil
		})
		btn.Call("addEventListener", "click", onClick)
		bar.Call("appendChild", btn)
	}

	start := time.Now()
	var frame js.Func
	frame = js.FuncOf(func(this js.Value, _ []js.Value) interface{} {
		t := time.Since(start).Seconds()
		switch mode {
		case "julia":
			drawJulia(pixels, pal, t)
			blit(ctx, img, buf, pixels)
		case "lissajous":
			drawLissajous(ctx, t)
		default:
			drawPlasma(pixels, pal, t)
			blit(ctx, img, buf, pixels)
		}
		js.Global().Call("requestAnimationFrame", frame)
		return nil
	})
	js.Global().Call("requestAnimationFrame", frame)

	select {} // sleep forever; the browser drives everything via callbacks
}

func blit(ctx, img, buf js.Value, px []byte) {
	js.CopyBytesToJS(buf, px)
	img.Get("data").Call("set", buf)
	ctx.Call("putImageData", img, 0, 0)
}

func drawPlasma(px []byte, pal [][3]uint8, t float64) {
	i := 0
	for y := 0; y < height; y++ {
		fy := float64(y)
		for x := 0; x < width; x++ {
			fx := float64(x)
			v := math.Sin(fx/18+t) +
				math.Sin(fy/9-t*1.4) +
				math.Sin((fx+fy)/26+t*0.8) +
				math.Sin(math.Hypot(fx-width/2, fy-height/2)/14-t)
			idx := int((v + 4) * 32)
			if idx < 0 {
				idx = 0
			}
			if idx > 255 {
				idx = 255
			}
			c := pal[idx]
			px[i], px[i+1], px[i+2] = c[0], c[1], c[2]
			px[i+3] = 255
			i += 4
		}
	}
}

func drawJulia(px []byte, pal [][3]uint8, t float64) {
	angle := t * 0.15
	cr := math.Cos(angle) * 0.7885
	ci := math.Sin(angle) * 0.7885

	i := 0
	for y := 0; y < height; y++ {
		zy0 := (float64(y) - height/2) / (height / 2.6)
		for x := 0; x < width; x++ {
			zx := (float64(x) - width/2) / (width / 3.2)
			zy := zy0
			n := 0
			var zx2, zy2 float64
			for ; n < maxIter; n++ {
				zx2, zy2 = zx*zx, zy*zy
				if zx2+zy2 > 4 {
					break
				}
				zy = 2*zx*zy + ci
				zx = zx2 - zy2 + cr
			}

			var c [3]uint8
			if n >= maxIter {
				c = [3]uint8{12, 9, 36}
			} else {
				mu := float64(n) + 1 - math.Log2(math.Log(math.Sqrt(zx2+zy2)))
				if mu < 0 {
					mu = 0
				}
				k := int(mu*16) % 256
				if k < 0 {
					k += 256
				}
				c = pal[k]
			}
			px[i], px[i+1], px[i+2] = c[0], c[1], c[2]
			px[i+3] = 255
			i += 4
		}
	}
}

func drawLissajous(ctx js.Value, t float64) {
	ctx.Set("fillStyle", "#0b0920")
	ctx.Call("fillRect", 0, 0, width, height)
	ctx.Set("strokeStyle", "#8f6bff")
	ctx.Set("lineWidth", 1.5)
	ctx.Call("beginPath")

	const steps = 1400
	for i := 0; i <= steps; i++ {
		u := 2 * math.Pi * float64(i) / steps
		x := width/2 + (width/2-6)*math.Sin(3*u+t*0.9)
		y := height/2 + (height/2-6)*math.Sin(4*u)
		if i == 0 {
			ctx.Call("moveTo", x, y)
		} else {
			ctx.Call("lineTo", x, y)
		}
	}
	ctx.Call("stroke")
}

// makePalette builds a bright HSV rainbow with n entries.
func makePalette(n int) [][3]uint8 {
	pal := make([][3]uint8, n)
	for i := 0; i < n; i++ {
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
		pal[i] = [3]uint8{r, g, b}
	}
	return pal
}
