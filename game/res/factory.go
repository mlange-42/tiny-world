package res

import (
	"image"
	"math"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/terr"
)

type EntityFactory struct {
	landUseBuilder    generic.Map4[comp.Tile, comp.Terrain, comp.UpdateTick, comp.RandomSprite]
	productionBuilder generic.Map6[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Production, comp.Consumption, comp.RandomSprite]
	warehouseBuilder  generic.Map5[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Warehouse, comp.RandomSprite]
	pathBuilder       generic.Map4[comp.Tile, comp.Terrain, comp.Path, comp.RandomSprite]

	terrain         generic.Resource[Terrain]
	terrainEntities generic.Resource[TerrainEntities]
	landUse         generic.Resource[LandUse]
	landUseEntities generic.Resource[LandUseEntities]

	update generic.Resource[UpdateInterval]
}

func NewEntityFactory(world *ecs.World) EntityFactory {
	return EntityFactory{
		landUseBuilder:    generic.NewMap4[comp.Tile, comp.Terrain, comp.UpdateTick, comp.RandomSprite](world),
		productionBuilder: generic.NewMap6[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Production, comp.Consumption, comp.RandomSprite](world),
		warehouseBuilder:  generic.NewMap5[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Warehouse, comp.RandomSprite](world),
		pathBuilder:       generic.NewMap4[comp.Tile, comp.Terrain, comp.Path, comp.RandomSprite](world),

		terrain:         generic.NewResource[Terrain](world),
		terrainEntities: generic.NewResource[TerrainEntities](world),
		landUse:         generic.NewResource[LandUse](world),
		landUseEntities: generic.NewResource[LandUseEntities](world),

		update: generic.NewResource[UpdateInterval](world),
	}
}

func (f *EntityFactory) createLandUse(pos image.Point, t terr.Terrain) ecs.Entity {
	e := f.landUseBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.UpdateTick{Tick: rand.Int63n(f.update.Get().Interval)},
		&comp.RandomSprite{Rand: uint16(rand.Int31n(math.MaxUint16))},
	)
	return e
}

func (f *EntityFactory) createWarehouse(pos image.Point, t terr.Terrain) ecs.Entity {
	e := f.warehouseBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.UpdateTick{Tick: rand.Int63n(f.update.Get().Interval)},
		&comp.Warehouse{},
		&comp.RandomSprite{Rand: uint16(rand.Int31n(math.MaxUint16))},
	)
	return e
}

func (f *EntityFactory) createPath(pos image.Point, t terr.Terrain) ecs.Entity {
	e := f.pathBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.Path{Haulers: []comp.HaulerEntry{}},
		&comp.RandomSprite{Rand: uint16(rand.Int31n(math.MaxUint16))},
	)
	return e
}

func (f *EntityFactory) createProduction(pos image.Point, t terr.Terrain, prod *terr.Production) ecs.Entity {
	update := f.update.Get()
	e := f.productionBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.UpdateTick{Tick: rand.Int63n(update.Interval)},
		&comp.Production{Type: prod.Resource, Amount: 0, Countdown: update.Countdown},
		&comp.Consumption{Amount: prod.FoodConsumption, Countdown: update.Countdown},
		&comp.RandomSprite{Rand: uint16(rand.Int31n(math.MaxUint16))},
	)
	return e
}

func (f *EntityFactory) Create(pos image.Point, t terr.Terrain) ecs.Entity {
	props := &terr.Properties[t]
	if props.IsWarehouse {
		return f.createWarehouse(pos, t)
	}
	if props.IsPath {
		return f.createPath(pos, t)
	}
	prod := props.Production
	if prod.MaxProduction == 0 {
		return f.createLandUse(pos, t)
	}
	return f.createProduction(pos, t, &prod)
}

func (f *EntityFactory) Set(world *ecs.World, x, y int, value terr.Terrain) {
	if !terr.Properties[value].IsTerrain {
		f.landUse.Get().Set(x, y, value)
		e := f.Create(image.Pt(x, y), value)
		f.landUseEntities.Get().Set(x, y, e)
		return
	}
	t := f.terrain.Get()
	tE := f.terrainEntities.Get()

	eHere := tE.Get(x, y)
	if !eHere.IsZero() {
		world.RemoveEntity(eHere)
	}

	t.Set(x, y, value)
	e := f.Create(image.Pt(x, y), value)
	tE.Set(x, y, e)

	f.setNeighbor(t, tE, x-1, y)
	f.setNeighbor(t, tE, x+1, y)
	f.setNeighbor(t, tE, x, y-1)
	f.setNeighbor(t, tE, x, y+1)
}

func (f *EntityFactory) setNeighbor(t *Terrain, tE *TerrainEntities, x, y int) {
	if t.Contains(x, y) && t.Get(x, y) == terr.Air {
		t.Set(x, y, terr.Buildable)
		e := f.Create(image.Pt(x, y), terr.Buildable)
		tE.Set(x, y, e)
	}
}
