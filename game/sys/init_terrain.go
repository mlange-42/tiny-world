package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game"
	"github.com/mlange-42/tiny-world/game/terr"
)

// InitTerrain system.
type InitTerrain struct {
}

// Initialize the system
func (s *InitTerrain) Initialize(world *ecs.World) {
	terrain := generic.NewResource[game.Terrain](world)

	t := terrain.Get()
	t.Set(t.Width()/2, t.Height()/2, terr.Grass)
}

// Update the system
func (s *InitTerrain) Update(world *ecs.World) {}

// Finalize the system
func (s *InitTerrain) Finalize(world *ecs.World) {}
