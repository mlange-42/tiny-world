package sys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

type mouse struct {
	X int
	Y int
}

// PanAndZoom system.
type PanAndZoom struct {
	PanButton ebiten.MouseButton

	mouseStart mouse

	view    generic.Resource[res.View]
	terrain generic.Resource[res.Terrain]
	sprites generic.Resource[res.Sprites]
}

// Initialize the system
func (s *PanAndZoom) Initialize(world *ecs.World) {
	s.view = generic.NewResource[res.View](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.sprites = generic.NewResource[res.Sprites](world)
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
	mx, my := view.ScreenToGlobal(x, y)
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
