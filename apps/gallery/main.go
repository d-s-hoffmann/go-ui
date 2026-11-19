//go:build js && wasm

// Command gallery shows animated scenes on a canvas in the browser.
package main

import (
	"image"
	"syscall/js"
	"time"

	"github.com/d-s-hoffmann/go-ui/internal/scene"
)

const (
	width  = 320
	height = 200
)

func main() {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "stage")
	if canvas.IsNull() {
		println("no #stage canvas found")
		return
	}
	ctx := canvas.Call("getContext", "2d")

	frame := image.NewRGBA(image.Rect(0, 0, width, height))
	img := ctx.Call("createImageData", width, height)
	buf := js.Global().Get("Uint8Array").New(len(frame.Pix))

	scenes := []scene.Scene{scene.Plasma{}, scene.Julia{}, &scene.Lissajous{}}
	current := scenes[0]

	bar := doc.Call("getElementById", "controls")
	for _, s := range scenes {
		btn := doc.Call("createElement", "button")
		btn.Set("textContent", s.Name())
		btn.Set("className", "btn")
		onClick := js.FuncOf(func(this js.Value, _ []js.Value) any {
			current = s
			return nil
		})
		btn.Call("addEventListener", "click", onClick)
		bar.Call("appendChild", btn)
	}

	start := time.Now()
	var tick js.Func
	tick = js.FuncOf(func(this js.Value, _ []js.Value) any {
		current.Render(frame, time.Since(start))
		blit(ctx, img, buf, frame)
		js.Global().Call("requestAnimationFrame", tick)
		return nil
	})
	js.Global().Call("requestAnimationFrame", tick)

	select {} // sleep forever; the browser drives everything via callbacks
}

// blit copies frame onto the canvas through the reusable ImageData img
// and its staging buffer buf.
func blit(ctx, img, buf js.Value, frame *image.RGBA) {
	js.CopyBytesToJS(buf, frame.Pix)
	img.Get("data").Call("set", buf)
	ctx.Call("putImageData", img, 0, 0)
}
