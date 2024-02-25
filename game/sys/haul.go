package sys

import (
	ares "github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/nav"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Haul system.
type Haul struct {
	time     generic.Resource[ares.Tick]
	update   generic.Resource[res.UpdateInterval]
	stock    generic.Resource[res.Stock]
	landUse  generic.Resource[res.LandUse]
	landUseE generic.Resource[res.LandUseEntities]

	filter    generic.Filter3[comp.Tile, comp.Hauler, comp.UpdateTick]
	haulerMap generic.Map2[comp.Tile, comp.Hauler]
	homeMap   generic.Map2[comp.Tile, comp.Production]

	aStar nav.AStar

	arrived []ecs.Entity
}

// Initialize the system
func (s *Haul) Initialize(world *ecs.World) {
	s.time = generic.NewResource[ares.Tick](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.landUseE = generic.NewResource[res.LandUseEntities](world)

	s.filter = *generic.NewFilter3[comp.Tile, comp.Hauler, comp.UpdateTick]()
	s.haulerMap = generic.NewMap2[comp.Tile, comp.Hauler](world)
	s.homeMap = generic.NewMap2[comp.Tile, comp.Production](world)

	s.aStar = nav.NewAStar(s.landUse.Get())
}

// Update the system
func (s *Haul) Update(world *ecs.World) {
	tick := s.time.Get().Tick
	update := s.update.Get()
	landUse := s.landUse.Get()
	stock := s.stock.Get()
	tickMod := tick % update.Interval

	query := s.filter.Query(world)
	for query.Next() {
		tile, haul, up := query.Get()

		if up.Tick != tickMod {
			continue
		}

		haul.Path = haul.Path[:len(haul.Path)-1]
		last := haul.Path[len(haul.Path)-1]
		tile.X, tile.Y = last.X, last.Y

		if len(haul.Path) <= 1 {
			s.arrived = append(s.arrived, query.Entity())
		}
	}

	for _, e := range s.arrived {
		tile, haul := s.haulerMap.Get(e)

		if !world.Alive(haul.Home) {
			world.RemoveEntity(e)
			continue
		}

		home, prod := s.homeMap.Get(haul.Home)
		if landUse.Get(tile.X, tile.Y) == terr.Warehouse {
			stock.Res[haul.Hauls]++

			path, ok := s.aStar.FindPath(*tile, *home)
			if !ok {
				prod.Paused = false
				world.RemoveEntity(e)
			}
			haul.Path = path
			continue
		}

		prod.Paused = false
		world.RemoveEntity(e)
	}

	s.arrived = s.arrived[:0]
}

// Finalize the system
func (s *Haul) Finalize(world *ecs.World) {}
