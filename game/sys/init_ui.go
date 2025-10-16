package sys

import (
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/res"
)

// InitUI system.
type InitUI struct {
	ui res.UI
}

// Initialize the system
func (s *InitUI) Initialize(world *ecs.World) {
	s.ui = res.NewUI(world,
		ecs.GetResource[res.Selection](world),
		ecs.GetResource[res.Fonts](world),
		ecs.GetResource[res.Sprites](world),
		ecs.GetResource[res.RandomTerrains](world),
		ecs.GetResource[res.SaveEvent](world),
		ecs.GetResource[res.EditorMode](world))

	ecs.AddResource(world, &s.ui)
}

// Update the system
func (s *InitUI) Update(world *ecs.World) {}

// Finalize the system
func (s *InitUI) Finalize(world *ecs.World) {}
