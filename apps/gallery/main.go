//go:build js && wasm

// Command gallery shows animated scenes on a canvas in the browser.
//
// It only wires packages together: scenes come from internal/scene, and
// everything that talks to the browser lives in internal/web.
package main

import (
	"time"

	"github.com/d-s-hoffmann/go-ui/internal/scene"
	"github.com/d-s-hoffmann/go-ui/internal/web"
)

func main() {
	canvas, err := web.NewCanvas("stage")
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
	if err := web.Buttons("controls", names, func(i int) { current = scenes[i] }); err != nil {
		println(err.Error())
		return
	}

	web.Animate(func(t time.Duration) {
		current.Render(canvas.Frame(), t)
		canvas.Present()
	})

	select {} // keep the program alive; the browser drives it via callbacks
}
