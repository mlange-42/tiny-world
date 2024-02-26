package sys

import (
	"math"

	ares "github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/nav"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
)

// DoProduction system.
type DoProduction struct {
	time    generic.Resource[ares.Tick]
	update  generic.Resource[res.UpdateInterval]
	stock   generic.Resource[res.Stock]
	landUse generic.Resource[res.LandUse]

	filter          generic.Filter3[comp.Tile, comp.UpdateTick, comp.Production]
	warehouseFilter generic.Filter1[comp.Tile]
	markerBuilder   generic.Map2[comp.Tile, comp.ProductionMarker]
	haulerBuilder   generic.Map2[comp.Tile, comp.Hauler]
	productionMap   generic.Map1[comp.Production]

	aStar nav.AStar

	warehouses []comp.Tile
	toCreate   []markerEntry
}

// Initialize the system
func (s *DoProduction) Initialize(world *ecs.World) {
	s.time = generic.NewResource[ares.Tick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.landUse = generic.NewResource[res.LandUse](world)

	s.filter = *generic.NewFilter3[comp.Tile, comp.UpdateTick, comp.Production]()
	s.warehouseFilter = *generic.NewFilter1[comp.Tile]().With(generic.T[comp.Warehouse]())
	s.markerBuilder = generic.NewMap2[comp.Tile, comp.ProductionMarker](world)
	s.haulerBuilder = generic.NewMap2[comp.Tile, comp.Hauler](world)
	s.productionMap = generic.NewMap1[comp.Production](world)

	s.aStar = nav.NewAStar(s.landUse.Get())
}

// Update the system
func (s *DoProduction) Update(world *ecs.World) {
	//stock := s.stock.Get()
	tick := s.time.Get().Tick
	update := s.update.Get()
	tickMod := tick % update.Interval

	query := s.filter.Query(world)
	for query.Next() {
		tile, up, pr := query.Get()

		if up.Tick != tickMod {
			continue
		}

		if pr.Paused {
			continue
		}

		pr.Countdown -= pr.Amount
		if pr.Countdown < 0 {
			pr.Countdown += update.Countdown
			//stock.Res[pr.Type]++
			s.toCreate = append(s.toCreate, markerEntry{Tile: *tile, Resource: pr.Type, Home: query.Entity()})
		}
	}

	if len(s.toCreate) > 0 {
		query := s.warehouseFilter.Query(world)
		for query.Next() {
			s.warehouses = append(s.warehouses, *query.Get())
		}
	}

	for _, entry := range s.toCreate {
		s.markerBuilder.NewWith(
			&entry.Tile,
			&comp.ProductionMarker{StartTick: tick, Resource: entry.Resource},
		)

		var bestPath []comp.Tile
		bestPathLen := math.MaxInt
		for _, tile := range s.warehouses {
			if path, ok := s.aStar.FindPath(entry.Tile, tile); ok {
				if len(path) < bestPathLen {
					bestPathLen = len(path)
					bestPath = path
				}
			}
		}
		if len(bestPath) == 0 {
			continue
		}
		prod := s.productionMap.Get(entry.Home)
		prod.Paused = true
		s.haulerBuilder.NewWith(
			&entry.Tile,
			&comp.Hauler{
				Hauls:        entry.Resource,
				Home:         entry.Home,
				Path:         bestPath,
				PathFraction: uint8(update.Interval / 2),
			},
		)
	}
	s.warehouses = s.warehouses[:0]
	s.toCreate = s.toCreate[:0]
}

// Finalize the system
func (s *DoProduction) Finalize(world *ecs.World) {}

type markerEntry struct {
	Tile     comp.Tile
	Resource resource.Resource
	Home     ecs.Entity
}
