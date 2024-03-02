package menu

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// UpdateUI system.
type UpdateUI struct {
	ui generic.Resource[UI]
}

// Initialize the system
func (s *UpdateUI) Initialize(world *ecs.World) {
	s.ui = generic.NewResource[UI](world)
}

// Update the system
func (s *UpdateUI) Update(world *ecs.World) {
	ui := s.ui.Get()

	ui.UI().Update()
}

// Finalize the system
func (s *UpdateUI) Finalize(world *ecs.World) {}
