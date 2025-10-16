package sys

import (
	"math"

	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/nav"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/sprites"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Haul system.
type Haul struct {
	speed    ecs.Resource[res.GameSpeed]
	update   ecs.Resource[res.UpdateInterval]
	stock    ecs.Resource[res.Stock]
	landUse  ecs.Resource[res.LandUse]
	landUseE ecs.Resource[res.LandUseEntities]

	prodFilter      *ecs.Filter3[comp.Tile, comp.Terrain, comp.Production]
	warehouseFilter *ecs.Filter2[comp.Tile, comp.Terrain]
	filter          *ecs.Filter2[comp.Tile, comp.Hauler]

	haulerMap     *ecs.Map2[comp.Tile, comp.Hauler]
	homeMap       *ecs.Map3[comp.Tile, comp.Terrain, comp.Production]
	haulerBuilder *ecs.Map3[comp.Tile, comp.Hauler, comp.HaulerSprite]
	productionMap *ecs.Map2[comp.Terrain, comp.Production]

	aStar nav.AStar

	warehouses [][]comp.Tile
	toCreate   []markerEntry
	arrived    []ecs.Entity

	haulerSprites []int
}

// Initialize the system
func (s *Haul) Initialize(world *ecs.World) {
	s.speed = ecs.NewResource[res.GameSpeed](world)
	s.update = ecs.NewResource[res.UpdateInterval](world)
	s.stock = ecs.NewResource[res.Stock](world)
	s.landUse = ecs.NewResource[res.LandUse](world)
	s.landUseE = ecs.NewResource[res.LandUseEntities](world)

	s.prodFilter = ecs.NewFilter3[comp.Tile, comp.Terrain, comp.Production](world)
	s.warehouseFilter = ecs.NewFilter2[comp.Tile, comp.Terrain](world).With(ecs.C[comp.Warehouse]())
	s.filter = ecs.NewFilter2[comp.Tile, comp.Hauler](world)

	s.haulerMap = ecs.NewMap2[comp.Tile, comp.Hauler](world)
	s.homeMap = ecs.NewMap3[comp.Tile, comp.Terrain, comp.Production](world)
	s.haulerBuilder = ecs.NewMap3[comp.Tile, comp.Hauler, comp.HaulerSprite](world)
	s.productionMap = ecs.NewMap2[comp.Terrain, comp.Production](world)

	s.aStar = nav.NewAStar(s.landUse.Get())

	spritesRes := ecs.NewResource[res.Sprites](world)
	spr := spritesRes.Get()

	s.haulerSprites = make([]int, len(terr.Properties))
	for i := range terr.Properties {
		s.haulerSprites[i] = spr.GetIndex(sprites.HaulerPrefix + terr.Properties[i].Name)
	}

	s.warehouses = make([][]comp.Tile, len(resource.Properties))
}

// Update the system
func (s *Haul) Update(world *ecs.World) {
	if s.speed.Get().Pause {
		return
	}

	update := s.update.Get()
	landUse := s.landUse.Get()
	landUseE := s.landUseE.Get()
	stock := s.stock.Get()

	prodQuery := s.prodFilter.Query()
	for prodQuery.Next() {
		tile, tp, prod := prodQuery.Get()
		if prod.Stock < terr.Properties[tp.Terrain].Production.HaulCapacity || prod.IsHauling {
			continue
		}
		s.toCreate = append(s.toCreate, markerEntry{Tile: *tile, Resource: prod.Resource, Home: prodQuery.Entity()})
	}

	query := s.filter.Query()
	for query.Next() {
		tile, haul := query.Get()

		haul.PathFraction++
		if haul.Index <= 1 && haul.PathFraction >= uint8(update.Interval-1) {
			s.arrived = append(s.arrived, query.Entity())
			continue
		}

		if haul.PathFraction < uint8(update.Interval) {
			continue
		}
		haul.PathFraction = 0

		haul.Index--
		last := haul.Path[haul.Index]
		tile.X, tile.Y = last.X, last.Y
	}

	if len(s.toCreate) > 0 {
		query := s.warehouseFilter.Query()
		for query.Next() {
			tile, ter := query.Get()
			storage := terr.Properties[ter.Terrain].Storage
			for i, st := range storage {
				if st > 0 {
					s.warehouses[i] = append(s.warehouses[i], *tile)
				}
			}
		}
	}

	for _, entry := range s.toCreate {
		var bestPath []comp.Tile
		bestPathLen := math.MaxInt
		for _, tile := range s.warehouses[entry.Resource] {
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
		luHere := landUse.Get(entry.Tile.X, entry.Tile.Y)

		tp, prod := s.productionMap.Get(entry.Home)
		prod.Stock -= terr.Properties[tp.Terrain].Production.HaulCapacity
		prod.IsHauling = true
		s.haulerBuilder.NewEntity(
			&entry.Tile,
			&comp.Hauler{
				Hauls:        entry.Resource,
				Home:         entry.Home,
				Path:         bestPath,
				PathFraction: 0,
				Index:        len(bestPath) - 1,
			},
			&comp.HaulerSprite{
				SpriteIndex: s.haulerSprites[luHere],
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

		home, tp, prod := s.homeMap.Get(haul.Home)
		if terr.Properties[landUse.Get(target.X, target.Y)].TerrainBits.Contains(terr.IsWarehouse) {
			amount := int(terr.Properties[tp.Terrain].Production.HaulCapacity)
			stock.AddResources(haul.Hauls, amount)

			path, ok := s.aStar.FindPath(target, *home)
			if !ok {
				prod.IsHauling = false
				world.RemoveEntity(e)
			}
			haul.Path = path
			haul.Index = len(path) - 1
			haul.PathFraction = uint8(update.Interval/2) + 1
			*tile = target

			if haul.Home != landUseE.Get(target.X, target.Y) {
				continue
			}
		}

		prod.IsHauling = false
		world.RemoveEntity(e)
	}

	for i := range s.warehouses {
		s.warehouses[i] = s.warehouses[i][:0]
	}
	s.toCreate = s.toCreate[:0]
	s.arrived = s.arrived[:0]
}

// Finalize the system
func (s *Haul) Finalize(world *ecs.World) {}
