package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// InitTerrain system.
type InitTerrain struct {
}

// Initialize the system
func (s *InitTerrain) Initialize(world *ecs.World) {
	factory := generic.NewResource[res.EntityFactory](world)
	terrain := generic.NewResource[res.Terrain](world)

	t := terrain.Get()
	x, y := t.Width()/2, t.Height()/2

	fac := factory.Get()

	fac.Set(world, x, y, terr.Default)
	fac.Set(world, x, y, terr.Warehouse)
}

// Update the system
func (s *InitTerrain) Update(world *ecs.World) {}

// Finalize the system
func (s *InitTerrain) Finalize(world *ecs.World) {}
