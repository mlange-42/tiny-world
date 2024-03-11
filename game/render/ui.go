package render

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// UI is a system to render the user interface.
type UI struct {
	screen generic.Resource[res.Screen]
	ui     generic.Resource[res.UI]
}

// InitializeUI the system
func (s *UI) InitializeUI(world *ecs.World) {
	s.ui = generic.NewResource[res.UI](world)
	s.screen = generic.NewResource[res.Screen](world)
}

// UpdateUI the system
func (s *UI) UpdateUI(world *ecs.World) {
	screen := s.screen.Get()
	ui := s.ui.Get()

	ui.UI().Draw(screen.Image)
}

// PostUpdateUI the system
func (s *UI) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *UI) FinalizeUI(world *ecs.World) {}
