package sys

import (
	"image"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/math"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/util"
)

type mouse struct {
	X int
	Y int
}

// PanAndZoom system.
type PanAndZoom struct {
	PanButton        ebiten.MouseButton
	KeyboardPanSpeed int

	MinZoom float64
	MaxZoom float64

	mouseStart mouse

	view    generic.Resource[res.View]
	screen  generic.Resource[res.Screen]
	bounds  generic.Resource[res.WorldBounds]
	terrain generic.Resource[res.Terrain]
	sprites generic.Resource[res.Sprites]

	inputChars []rune
}

// Initialize the system
func (s *PanAndZoom) Initialize(world *ecs.World) {
	s.view = generic.NewResource[res.View](world)
	s.screen = generic.NewResource[res.Screen](world)
	s.bounds = generic.NewResource[res.WorldBounds](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.sprites = generic.NewResource[res.Sprites](world)
}

// Update the system
func (s *PanAndZoom) Update(world *ecs.World) {
	view := s.view.Get()
	screen := s.screen.Get()
	bounds := s.bounds.Get()

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

	panSpeed := math.MaxInt(int(float64(s.KeyboardPanSpeed)/view.Zoom), 1)
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		view.X += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		view.X -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		view.Y -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) ||
		(ebiten.IsKeyPressed(ebiten.KeyS) && !ebiten.IsKeyPressed(ebiten.KeyControl)) {
		view.Y += panSpeed
	}

	_, dy := ebiten.Wheel()
	x, y := ebiten.CursorPosition()
	mx, my := view.ScreenToGlobal(x, y)

	s.inputChars = ebiten.AppendInputChars(s.inputChars)
	if (dy > 0 || slices.Contains(s.inputChars, '+')) && view.Zoom < s.MaxZoom {
		view.Zoom *= 2
		view.X += (mx - view.X) / 2
		view.Y += (my - view.Y) / 2
	}
	if (dy < 0 || slices.Contains(s.inputChars, '-')) && view.Zoom > s.MinZoom {
		view.Zoom /= 2
		view.X -= (mx - view.X)
		view.Y -= (my - view.Y)
	}
	s.inputChars = s.inputChars[:0]

	glBounds := view.BoundsToGlobal(bounds)
	center := image.Pt(view.ScreenToGlobal(screen.Width/2, screen.Height/2))
	if !center.In(glBounds) {
		center := util.Clamp(glBounds, center)
		view.X = center.X - int(float64(screen.Width/2)/view.Zoom)
		view.Y = center.Y - int(float64(screen.Height/2)/view.Zoom)
	}
}

// Finalize the system
func (s *PanAndZoom) Finalize(world *ecs.World) {}
