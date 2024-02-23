package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// UpdateUI system.
type UpdateUI struct {
	hud generic.Resource[res.HUD]
	ui  generic.Resource[res.UI]
}

// Initialize the system
func (s *UpdateUI) Initialize(world *ecs.World) {
	s.hud = generic.NewResource[res.HUD](world)
	s.ui = generic.NewResource[res.UI](world)
}

// Update the system
func (s *UpdateUI) Update(world *ecs.World) {
	hud := s.hud.Get()
	ui := s.ui.Get()

	ui.UI.Update()
	hud.UI.Update()
}

// Finalize the system
func (s *UpdateUI) Finalize(world *ecs.World) {}
