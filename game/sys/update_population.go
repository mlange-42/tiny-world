package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/math"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// UpdateProduction system.
type UpdatePopulation struct {
	time    generic.Resource[res.GameTick]
	speed   generic.Resource[res.GameSpeed]
	update  generic.Resource[res.UpdateInterval]
	terrain generic.Resource[res.Terrain]
	landUse generic.Resource[res.LandUse]

	filter generic.Filter3[comp.Tile, comp.UpdateTick, comp.PopulationSupport]
}

// Initialize the system
func (s *UpdatePopulation) Initialize(world *ecs.World) {
	s.time = generic.NewResource[res.GameTick](world)
	s.speed = generic.NewResource[res.GameSpeed](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.landUse = generic.NewResource[res.LandUse](world)

	s.filter = *generic.NewFilter3[comp.Tile, comp.UpdateTick, comp.PopulationSupport]()
}

// Update the system
func (s *UpdatePopulation) Update(world *ecs.World) {
	if s.speed.Get().Pause {
		return
	}

	terrain := s.terrain.Get()
	landUse := s.landUse.Get()
	tick := s.time.Get().Tick
	interval := s.update.Get().Interval
	tickMod := tick % interval

	query := s.filter.Query(world)
	for query.Next() {
		tile, up, pop := query.Get()

		if up.Tick != tickMod {
			continue
		}
		pop.Pop = 0

		lu := landUse.Get(tile.X, tile.Y)

		supp := &terr.Properties[lu].PopulationSupport
		if supp.RequiredTerrain != terr.Air &&
			terrain.CountNeighbors4(tile.X, tile.Y, supp.RequiredTerrain) == 0 &&
			landUse.CountNeighbors4(tile.X, tile.Y, supp.RequiredTerrain) == 0 {
			pop.HasRequired = false
			continue
		}
		pop.HasRequired = true
		count := int(supp.BasePopulation)
		if supp.BonusTerrain != 0 {
			count += terrain.CountNeighborsMask8(tile.X, tile.Y, supp.BonusTerrain) +
				landUse.CountNeighborsMask8(tile.X, tile.Y, supp.BonusTerrain)
		}
		if supp.MalusTerrain != 0 {
			count -= terrain.CountNeighborsMask8(tile.X, tile.Y, supp.MalusTerrain) +
				landUse.CountNeighborsMask8(tile.X, tile.Y, supp.MalusTerrain)
		}
		pop.Pop = uint8(math.ClampInt(count, 0, int(supp.MaxPopulation)))
	}
}

// Finalize the system
func (s *UpdatePopulation) Finalize(world *ecs.World) {}
