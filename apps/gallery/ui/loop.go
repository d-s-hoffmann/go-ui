//go:build js && wasm

package ui

import (
	"syscall/js"
	"time"
)

// Animate calls frame once per display refresh, via requestAnimationFrame,
// with the time elapsed since Animate was called. It returns immediately;
// the browser drives the loop from then on.
func Animate(frame func(t time.Duration)) {
	start := time.Now()
	var tick js.Func
	tick = js.FuncOf(func(js.Value, []js.Value) any {
		frame(time.Since(start))
		js.Global().Call("requestAnimationFrame", tick)
		return nil
	})
	js.Global().Call("requestAnimationFrame", tick)
}
