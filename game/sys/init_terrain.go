package sys

import (
	"image"

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
	terrain := generic.NewResource[res.Terrain](world)
	factory := generic.NewResource[res.EntityFactory](world)
	landUse := generic.NewResource[res.LandUse](world)
	landUseE := generic.NewResource[res.LandUseEntities](world)

	t := terrain.Get()
	lu := landUse.Get()
	lue := landUseE.Get()
	x, y := t.Width()/2, t.Height()/2

	t.Set(x, y, terr.Grass)
	e := factory.Get().Create(image.Pt(x, y), terr.Warehouse)
	lu.Set(x, y, terr.Warehouse)
	lue.Set(x, y, e)
}

// Update the system
func (s *InitTerrain) Update(world *ecs.World) {}

// Finalize the system
func (s *InitTerrain) Finalize(world *ecs.World) {}
