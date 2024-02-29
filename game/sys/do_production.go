package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

// DoProduction system.
type DoProduction struct {
	speed   generic.Resource[res.GameSpeed]
	time    generic.Resource[res.GameTick]
	update  generic.Resource[res.UpdateInterval]
	stock   generic.Resource[res.Stock]
	landUse generic.Resource[res.LandUse]

	filter        generic.Filter4[comp.Terrain, comp.Tile, comp.UpdateTick, comp.Production]
	markerBuilder generic.Map2[comp.Tile, comp.ProductionMarker]

	toCreate []markerEntry
}

// Initialize the system
func (s *DoProduction) Initialize(world *ecs.World) {
	s.speed = generic.NewResource[res.GameSpeed](world)
	s.time = generic.NewResource[res.GameTick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.landUse = generic.NewResource[res.LandUse](world)

	s.filter = *generic.NewFilter4[comp.Terrain, comp.Tile, comp.UpdateTick, comp.Production]()
	s.markerBuilder = generic.NewMap2[comp.Tile, comp.ProductionMarker](world)
}

// Update the system
func (s *DoProduction) Update(world *ecs.World) {
	if s.speed.Get().Pause {
		return
	}

	tick := s.time.Get().Tick
	update := s.update.Get()
	tickMod := tick % update.Interval

	query := s.filter.Query(world)
	for query.Next() {
		ter, tile, up, pr := query.Get()

		if up.Tick != tickMod {
			continue
		}

		if pr.Stock >= terr.Properties[ter.Terrain].Storage[pr.Resource] {
			continue
		}

		pr.Countdown -= pr.Amount
		if pr.Countdown < 0 {
			pr.Countdown += update.Countdown
			pr.Stock++
			s.toCreate = append(s.toCreate, markerEntry{Tile: *tile, Resource: pr.Resource, Home: query.Entity()})
		}
	}

	for _, entry := range s.toCreate {
		s.markerBuilder.NewWith(
			&entry.Tile,
			&comp.ProductionMarker{StartTick: tick, Resource: entry.Resource},
		)
	}
	s.toCreate = s.toCreate[:0]
}

// Finalize the system
func (s *DoProduction) Finalize(world *ecs.World) {}

type markerEntry struct {
	Tile     comp.Tile
	Resource resource.Resource
	Home     ecs.Entity
}
