package res

import "github.com/hajimehoshi/ebiten/v2"

// EbitenImage resource for drawing.
type EbitenImage struct {
	Image  *ebiten.Image
	Width  int
	Height int
}
