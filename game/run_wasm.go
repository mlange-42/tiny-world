//go:build js

package game

import (
	"syscall/js"

	"github.com/mlange-42/ark-tools/app"
)

type canvasHelper struct {
	doc         js.Value
	canvas      js.Value
	mouseInside bool
}

func newCanvasHelper() *canvasHelper {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementsByTagName", "canvas").Index(0)

	helper := canvasHelper{
		doc:    doc,
		canvas: canvas,
	}

	canvas.Set("onmouseleave", js.FuncOf(helper.onMouseLeave))
	canvas.Set("onmouseenter", js.FuncOf(helper.onMouseEnter))

	return &helper
}

func (c *canvasHelper) isMouseInside(width, height int) bool {
	_, _ = width, height
	return c.mouseInside
}

func (c *canvasHelper) onMouseEnter(this js.Value, args []js.Value) interface{} {
	c.mouseInside = true
	return nil
}

func (c *canvasHelper) onMouseLeave(this js.Value, args []js.Value) interface{} {
	c.mouseInside = false
	return nil
}

func addRepl(app *app.App) {}
