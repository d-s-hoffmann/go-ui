//go:build js && wasm

// Command gallery shows animated scenes on a canvas in the browser.
//
// It only wires packages together: scenes come from ./scene, and
// everything that talks to the browser lives in ./ui.
package main

import (
	"time"

	"github.com/d-s-hoffmann/go-ui/apps/gallery/scene"
	"github.com/d-s-hoffmann/go-ui/apps/gallery/ui"
)

func main() {
	canvas, err := ui.NewCanvas("stage")
	if err != nil {
		println(err.Error())
		return
	}

	scenes := []scene.Scene{scene.Plasma{}, scene.Julia{}, &scene.Lissajous{}}
	current := scenes[0]

	names := make([]string, len(scenes))
	for i, s := range scenes {
		names[i] = s.Name()
	}
	if err := ui.Buttons("controls", names, func(i int) { current = scenes[i] }); err != nil {
		println(err.Error())
		return
	}

	ui.Animate(func(t time.Duration) {
		current.Render(canvas.Frame(), t)
		canvas.Present()
	})

	select {} // keep the program alive; the browser drives it via callbacks
}
