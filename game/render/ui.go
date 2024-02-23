package render

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// UI is a system to render the user interface.
type UI struct {
	screen generic.Resource[res.EbitenImage]
	hud    generic.Resource[res.HUD]
	ui     generic.Resource[res.UI]
}

// InitializeUI the system
func (s *UI) InitializeUI(world *ecs.World) {
	s.hud = generic.NewResource[res.HUD](world)
	s.ui = generic.NewResource[res.UI](world)
	s.screen = generic.NewResource[res.EbitenImage](world)
}

// UpdateUI the system
func (s *UI) UpdateUI(world *ecs.World) {
	screen := s.screen.Get()
	hud := s.hud.Get()
	ui := s.ui.Get()

	ui.UI.Draw(screen.Image)
	hud.UI.Draw(screen.Image)
}

// PostUpdateUI the system
func (s *UI) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *UI) FinalizeUI(world *ecs.World) {}
