package sys

import (
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/math"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// UpdateProduction system.
type UpdateProduction struct {
	time    ecs.Resource[res.GameTick]
	speed   ecs.Resource[res.GameSpeed]
	update  ecs.Resource[res.UpdateInterval]
	terrain ecs.Resource[res.Terrain]
	landUse ecs.Resource[res.LandUse]

	filter            *ecs.Filter4[comp.Tile, comp.UpdateTick, comp.Production, comp.Consumption]
	consumptionMapper *ecs.Map1[comp.Consumption]
}

// Initialize the system
func (s *UpdateProduction) Initialize(world *ecs.World) {
	s.time = ecs.NewResource[res.GameTick](world)
	s.speed = ecs.NewResource[res.GameSpeed](world)
	s.update = ecs.NewResource[res.UpdateInterval](world)
	s.terrain = ecs.NewResource[res.Terrain](world)
	s.landUse = ecs.NewResource[res.LandUse](world)

	s.filter = ecs.NewFilter4[comp.Tile, comp.UpdateTick, comp.Production, comp.Consumption](world)
	s.consumptionMapper = s.consumptionMapper.New(world)
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

	query := s.filter.Query()
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
			pr.HasRequired = false
			continue
		}
		pr.HasRequired = true
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
