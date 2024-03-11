package render

import (
	"image"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// UI is a system to render the user interface.
type CenterView struct {
	view    generic.Resource[res.View]
	screen  generic.Resource[res.Screen]
	terrain generic.Resource[res.Terrain]

	isInitialized bool
}

// InitializeUI the system
func (s *CenterView) InitializeUI(world *ecs.World) {
	s.view = generic.NewResource[res.View](world)
	s.screen = generic.NewResource[res.Screen](world)
	s.terrain = generic.NewResource[res.Terrain](world)
}

// UpdateUI the system
func (s *CenterView) UpdateUI(world *ecs.World) {
	if !s.isInitialized {
		view := s.view.Get()
		screen := s.screen.Get()
		terrain := s.terrain.Get()
		view.Center(image.Point{terrain.Width() / 2, terrain.Height() / 2}, screen.Width, screen.Height)
		s.isInitialized = true
	}
}

// PostUpdateUI the system
func (s *CenterView) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *CenterView) FinalizeUI(world *ecs.World) {}
