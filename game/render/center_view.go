package render

import (
	"image"

	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/res"
)

// CenterView is a system to render the user interface.
type CenterView struct {
	view    ecs.Resource[res.View]
	screen  ecs.Resource[res.Screen]
	terrain ecs.Resource[res.Terrain]

	isInitialized bool
}

// InitializeUI the system
func (s *CenterView) InitializeUI(world *ecs.World) {
	s.view = s.view.New(world)
	s.screen = s.screen.New(world)
	s.terrain = s.terrain.New(world)
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
