package sys

import (
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/res"
)

// UpdateUI system.
type UpdateUI struct {
	rules ecs.Resource[res.Rules]
	ui    ecs.Resource[res.UI]
}

// Initialize the system
func (s *UpdateUI) Initialize(world *ecs.World) {
	s.rules = ecs.NewResource[res.Rules](world)
	s.ui = ecs.NewResource[res.UI](world)

	rules := s.rules.Get()
	ui := s.ui.Get()
	ui.CreateRandomButtons(rules.RandomTerrainsCount)
}

// Update the system
func (s *UpdateUI) Update(world *ecs.World) {
	ui := s.ui.Get()

	ui.Update()
}

// Finalize the system
func (s *UpdateUI) Finalize(world *ecs.World) {}
