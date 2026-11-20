//go:build js && wasm

// Package ui connects Go code to the browser: it presents frames on a
// canvas, drives the animation loop and builds controls. It knows nothing
// about what is drawn.
package ui

import (
	"errors"
	"image"
	"strconv"
	"syscall/js"
)

// Canvas presents *image.RGBA frames on an HTML canvas element.
type Canvas struct {
	ctx   js.Value // CanvasRenderingContext2D
	img   js.Value // ImageData reused for every frame
	data  js.Value // img.data, the Uint8ClampedArray holding its pixels
	frame *image.RGBA
}

// NewCanvas returns a Canvas for the <canvas> element with the given id,
// sized to the element's width and height attributes.
func NewCanvas(id string) (*Canvas, error) {
	el := js.Global().Get("document").Call("getElementById", id)
	if el.IsNull() || el.Get("getContext").IsUndefined() {
		return nil, errors.New("ui: no canvas with id " + strconv.Quote(id))
	}
	w, h := el.Get("width").Int(), el.Get("height").Int()
	ctx := el.Call("getContext", "2d")
	img := ctx.Call("createImageData", w, h)
	return &Canvas{
		ctx:   ctx,
		img:   img,
		data:  img.Get("data"),
		frame: image.NewRGBA(image.Rect(0, 0, w, h)),
	}, nil
}

// Frame returns the image to draw the next frame into.
func (c *Canvas) Frame() *image.RGBA { return c.frame }

// Present copies the frame onto the canvas.
func (c *Canvas) Present() {
	js.CopyBytesToJS(c.data, c.frame.Pix)
	c.ctx.Call("putImageData", c.img, 0, 0)
}
