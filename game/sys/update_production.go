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
type UpdateProduction struct {
	time    generic.Resource[res.GameTick]
	speed   generic.Resource[res.GameSpeed]
	update  generic.Resource[res.UpdateInterval]
	terrain generic.Resource[res.Terrain]
	landUse generic.Resource[res.LandUse]

	filter generic.Filter4[comp.Tile, comp.UpdateTick, comp.Production, comp.Consumption]
}

// Initialize the system
func (s *UpdateProduction) Initialize(world *ecs.World) {
	s.time = generic.NewResource[res.GameTick](world)
	s.speed = generic.NewResource[res.GameSpeed](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.landUse = generic.NewResource[res.LandUse](world)

	s.filter = *generic.NewFilter4[comp.Tile, comp.UpdateTick, comp.Production, comp.Consumption]().Optional(generic.T[comp.Consumption]())
}

// Update the system
func (s *UpdateProduction) Update(world *ecs.World) {
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
		tile, up, pr, cons := query.Get()

		if up.Tick != tickMod {
			continue
		}
		pr.Amount = 0

		if !cons.IsSatisfied {
			continue
		}

		lu := landUse.Get(tile.X, tile.Y)

		prod := &terr.Properties[lu].Production
		if prod.RequiredTerrain != terr.Air &&
			terrain.CountNeighbors4(tile.X, tile.Y, prod.RequiredTerrain) == 0 &&
			landUse.CountNeighbors4(tile.X, tile.Y, prod.RequiredTerrain) == 0 {
			continue
		}
		count := 0
		if prod.ProductionTerrain != 0 {
			count += terrain.CountNeighborsMask8(tile.X, tile.Y, prod.ProductionTerrain) +
				landUse.CountNeighborsMask8(tile.X, tile.Y, prod.ProductionTerrain)
		}
		pr.Amount = uint8(math.MinInt(count, int(prod.MaxProduction)))
	}
}

// Finalize the system
func (s *UpdateProduction) Finalize(world *ecs.World) {}
