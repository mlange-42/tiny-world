package res

import "github.com/hajimehoshi/ebiten/v2"

// Screen resource for drawing.
type Screen struct {
	// The screen image.
	Image *ebiten.Image
	// Current screen width.
	Width int
	// Current screen height.
	Height int
}
