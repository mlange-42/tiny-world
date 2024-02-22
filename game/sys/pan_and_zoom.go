package sys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game"
)

type mouse struct {
	X int
	Y int
}

// PanAndZoom system.
type PanAndZoom struct {
	PanButton ebiten.MouseButton

	mouseStart mouse

	view    generic.Resource[game.View]
	terrain generic.Resource[game.Terrain]
	sprites generic.Resource[game.Sprites]
}

// Initialize the system
func (s *PanAndZoom) Initialize(world *ecs.World) {
	s.view = generic.NewResource[game.View](world)
	s.terrain = generic.NewResource[game.Terrain](world)
	s.sprites = generic.NewResource[game.Sprites](world)
}

// Update the system
func (s *PanAndZoom) Update(world *ecs.World) {
	view := s.view.Get()

	if inpututil.IsMouseButtonJustPressed(s.PanButton) {
		s.mouseStart.X, s.mouseStart.Y = ebiten.CursorPosition()
		return
	}

	if ebiten.IsMouseButtonPressed(s.PanButton) {
		x, y := ebiten.CursorPosition()
		view.X -= int(float64(x-s.mouseStart.X) / view.Zoom)
		view.Y -= int(float64(y-s.mouseStart.Y) / view.Zoom)

		s.mouseStart.X, s.mouseStart.Y = x, y
	}

	_, dy := ebiten.Wheel()
	x, y := ebiten.CursorPosition()
	mx, my := view.MouseToLocal(x, y)
	if dy > 0 && view.Zoom < 4 {
		view.Zoom *= 2
		view.X += (mx - view.X) / 2
		view.Y += (my - view.Y) / 2
	}
	if dy < 0 && view.Zoom > 0.25 {
		view.Zoom /= 2
		view.X -= (mx - view.X)
		view.Y -= (my - view.Y)
	}
}

// Finalize the system
func (s *PanAndZoom) Finalize(world *ecs.World) {}
