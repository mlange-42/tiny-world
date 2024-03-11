//go:build !js

package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type canvasHelper struct{}

func newCanvasHelper() *canvasHelper {
	return &canvasHelper{}
}

func (c *canvasHelper) isMouseInside(width, height int) bool {
	x, y := ebiten.CursorPosition()
	return x >= 0 && y >= 0 && x < width && y < height
}
