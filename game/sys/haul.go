package sys

import (
	"math"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/nav"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Haul system.
type Haul struct {
	update   generic.Resource[res.UpdateInterval]
	stock    generic.Resource[res.Stock]
	landUse  generic.Resource[res.LandUse]
	landUseE generic.Resource[res.LandUseEntities]

	prodFilter      generic.Filter2[comp.Tile, comp.Production]
	warehouseFilter generic.Filter1[comp.Tile]
	filter          generic.Filter2[comp.Tile, comp.Hauler]

	haulerMap     generic.Map2[comp.Tile, comp.Hauler]
	homeMap       generic.Map2[comp.Tile, comp.Production]
	haulerBuilder generic.Map2[comp.Tile, comp.Hauler]
	productionMap generic.Map1[comp.Production]

	aStar nav.AStar

	warehouses []comp.Tile
	toCreate   []markerEntry
	arrived    []ecs.Entity
}

// Initialize the system
func (s *Haul) Initialize(world *ecs.World) {
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.landUseE = generic.NewResource[res.LandUseEntities](world)

	s.prodFilter = *generic.NewFilter2[comp.Tile, comp.Production]()
	s.warehouseFilter = *generic.NewFilter1[comp.Tile]().With(generic.T[comp.Warehouse]())
	s.filter = *generic.NewFilter2[comp.Tile, comp.Hauler]()

	s.haulerMap = generic.NewMap2[comp.Tile, comp.Hauler](world)
	s.homeMap = generic.NewMap2[comp.Tile, comp.Production](world)
	s.haulerBuilder = generic.NewMap2[comp.Tile, comp.Hauler](world)
	s.productionMap = generic.NewMap1[comp.Production](world)

	s.aStar = nav.NewAStar(s.landUse.Get())
}

// Update the system
func (s *Haul) Update(world *ecs.World) {
	update := s.update.Get()
	landUse := s.landUse.Get()
	stock := s.stock.Get()

	prodQuery := s.prodFilter.Query(world)
	for prodQuery.Next() {
		tile, prod := prodQuery.Get()
		if prod.Stock == 0 || prod.IsHauling {
			continue
		}
		s.toCreate = append(s.toCreate, markerEntry{Tile: *tile, Resource: prod.Type, Home: prodQuery.Entity()})
	}

	query := s.filter.Query(world)
	for query.Next() {
		tile, haul := query.Get()

		haul.PathFraction++
		if len(haul.Path) <= 2 && haul.PathFraction >= uint8(update.Interval/2) {
			s.arrived = append(s.arrived, query.Entity())
			continue
		}

		if haul.PathFraction < uint8(update.Interval) {
			continue
		}
		haul.PathFraction = 0

		haul.Path = haul.Path[:len(haul.Path)-1]
		last := haul.Path[len(haul.Path)-1]
		tile.X, tile.Y = last.X, last.Y

		if len(haul.Path) <= 1 || (len(haul.Path) <= 2 && haul.PathFraction >= uint8(update.Interval/2)) {
			s.arrived = append(s.arrived, query.Entity())
		}
	}

	if len(s.toCreate) > 0 {
		query := s.warehouseFilter.Query(world)
		for query.Next() {
			s.warehouses = append(s.warehouses, *query.Get())
		}
	}

	for _, entry := range s.toCreate {
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
		prod.Stock -= 1
		prod.IsHauling = true
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

	for _, e := range s.arrived {
		tile, haul := s.haulerMap.Get(e)

		if !world.Alive(haul.Home) {
			world.RemoveEntity(e)
			continue
		}
		target := haul.Path[0]

		home, prod := s.homeMap.Get(haul.Home)
		if landUse.Get(target.X, target.Y) == terr.Warehouse {
			stock.Res[haul.Hauls]++

			path, ok := s.aStar.FindPath(target, *home)
			if !ok {
				prod.IsHauling = false
				world.RemoveEntity(e)
			}
			haul.Path = path
			haul.PathFraction = uint8(update.Interval / 2)
			*tile = target
			continue
		}

		prod.IsHauling = false
		world.RemoveEntity(e)
	}

	s.warehouses = s.warehouses[:0]
	s.toCreate = s.toCreate[:0]
	s.arrived = s.arrived[:0]
}

// Finalize the system
func (s *Haul) Finalize(world *ecs.World) {}
