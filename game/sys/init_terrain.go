package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// InitTerrain system.
type InitTerrain struct {
}

// Initialize the system
func (s *InitTerrain) Initialize(world *ecs.World) {
	rules := ecs.GetResource[res.Rules](world)
	fac := ecs.GetResource[res.EntityFactory](world)
	t := ecs.GetResource[res.Terrain](world)

	radiusMapper := generic.NewMap1[comp.BuildRadius](world)

	x, y := t.Width()/2, t.Height()/2

	fac.Set(world, x, y, terr.Default, 0)

	warehouse := fac.Set(world, x, y, terr.FirstBuilding, 0)
	radiusMapper.Assign(warehouse, &comp.BuildRadius{Radius: uint8(rules.InitialBuildRadius)})

	fac.SetBuildable(x, y, rules.InitialBuildRadius, true)
}

// Update the system
func (s *InitTerrain) Update(world *ecs.World) {}

// Finalize the system
func (s *InitTerrain) Finalize(world *ecs.World) {}
