//go:build js && wasm

package ui

import (
	"errors"
	"strconv"
	"syscall/js"
)

// Buttons appends one button per label to the element with the given id.
// Clicking the i-th button calls onSelect(i).
func Buttons(id string, labels []string, onSelect func(i int)) error {
	doc := js.Global().Get("document")
	bar := doc.Call("getElementById", id)
	if bar.IsNull() {
		return errors.New("ui: no element with id " + strconv.Quote(id))
	}
	for i, label := range labels {
		btn := doc.Call("createElement", "button")
		btn.Set("textContent", label)
		btn.Set("className", "btn")
		btn.Call("addEventListener", "click", js.FuncOf(func(js.Value, []js.Value) any {
			onSelect(i)
			return nil
		}))
		bar.Call("appendChild", btn)
	}
	return nil
}
