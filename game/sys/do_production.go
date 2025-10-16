package sys

import (
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

// DoProduction system.
type DoProduction struct {
	speed   ecs.Resource[res.GameSpeed]
	time    ecs.Resource[res.GameTick]
	update  ecs.Resource[res.UpdateInterval]
	stock   ecs.Resource[res.Stock]
	landUse ecs.Resource[res.LandUse]
	editor  ecs.Resource[res.EditorMode]

	filter        *ecs.Filter4[comp.Terrain, comp.Tile, comp.UpdateTick, comp.Production]
	markerBuilder *ecs.Map2[comp.Tile, comp.ProductionMarker]

	toCreate []markerEntry
}

// Initialize the system
func (s *DoProduction) Initialize(world *ecs.World) {
	s.speed = ecs.NewResource[res.GameSpeed](world)
	s.time = ecs.NewResource[res.GameTick](world)
	s.update = ecs.NewResource[res.UpdateInterval](world)
	s.stock = ecs.NewResource[res.Stock](world)
	s.landUse = ecs.NewResource[res.LandUse](world)
	s.editor = ecs.NewResource[res.EditorMode](world)

	s.filter = ecs.NewFilter4[comp.Terrain, comp.Tile, comp.UpdateTick, comp.Production](world)
	s.markerBuilder = ecs.NewMap2[comp.Tile, comp.ProductionMarker](world)
}

// Update the system
func (s *DoProduction) Update(world *ecs.World) {
	if s.speed.Get().Pause || s.editor.Get().IsEditor {
		return
	}

	tick := s.time.Get().Tick
	update := s.update.Get()
	tickMod := tick % update.Interval

	query := s.filter.Query()
	for query.Next() {
		ter, tile, up, pr := query.Get()

		if up.Tick != tickMod {
			continue
		}

		if pr.Stock >= terr.Properties[ter.Terrain].Storage[pr.Resource] {
			continue
		}

		pr.Countdown -= int(pr.Amount)
		if pr.Countdown < 0 {
			pr.Countdown += update.Countdown
			pr.Stock++
			s.toCreate = append(s.toCreate, markerEntry{Tile: *tile, Resource: pr.Resource, Home: query.Entity()})
		}
	}

	for _, entry := range s.toCreate {
		s.markerBuilder.NewEntity(
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
