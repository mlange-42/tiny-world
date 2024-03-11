package res

import (
	"image"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/terr"
)

// EntityFactory is a helper to create game entities.
type EntityFactory struct {
	landUseBuilder    generic.Map4[comp.Tile, comp.Terrain, comp.UpdateTick, comp.RandomSprite]
	productionBuilder generic.Map5[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Production, comp.RandomSprite]
	warehouseBuilder  generic.Map5[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Warehouse, comp.RandomSprite]
	pathBuilder       generic.Map4[comp.Tile, comp.Terrain, comp.Path, comp.RandomSprite]

	radiusMapper            generic.Map1[comp.BuildRadius]
	consumptionMapper       generic.Map1[comp.Consumption]
	populationMapper        generic.Map1[comp.Population]
	populationSupportMapper generic.Map1[comp.PopulationSupport]

	terrain         generic.Resource[Terrain]
	terrainEntities generic.Resource[TerrainEntities]
	landUse         generic.Resource[LandUse]
	landUseEntities generic.Resource[LandUseEntities]
	buildable       generic.Resource[Buildable]
	bounds          generic.Resource[WorldBounds]

	update generic.Resource[UpdateInterval]
}

// NewEntityFactory creates a new EntityFactory for a given world.
func NewEntityFactory(world *ecs.World) EntityFactory {
	return EntityFactory{
		landUseBuilder:    generic.NewMap4[comp.Tile, comp.Terrain, comp.UpdateTick, comp.RandomSprite](world),
		productionBuilder: generic.NewMap5[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Production, comp.RandomSprite](world),
		warehouseBuilder:  generic.NewMap5[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Warehouse, comp.RandomSprite](world),
		pathBuilder:       generic.NewMap4[comp.Tile, comp.Terrain, comp.Path, comp.RandomSprite](world),

		radiusMapper:            generic.NewMap1[comp.BuildRadius](world),
		consumptionMapper:       generic.NewMap1[comp.Consumption](world),
		populationMapper:        generic.NewMap1[comp.Population](world),
		populationSupportMapper: generic.NewMap1[comp.PopulationSupport](world),

		terrain:         generic.NewResource[Terrain](world),
		terrainEntities: generic.NewResource[TerrainEntities](world),
		landUse:         generic.NewResource[LandUse](world),
		landUseEntities: generic.NewResource[LandUseEntities](world),
		buildable:       generic.NewResource[Buildable](world),
		bounds:          generic.NewResource[WorldBounds](world),

		update: generic.NewResource[UpdateInterval](world),
	}
}

func (f *EntityFactory) createLandUse(pos image.Point, t terr.Terrain, randSprite uint16) ecs.Entity {
	e := f.landUseBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.UpdateTick{Tick: rand.Int63n(f.update.Get().Interval)},
		&comp.RandomSprite{Rand: randSprite},
	)
	return e
}

func (f *EntityFactory) createWarehouse(pos image.Point, t terr.Terrain, randSprite uint16) ecs.Entity {
	e := f.warehouseBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.UpdateTick{Tick: rand.Int63n(f.update.Get().Interval)},
		&comp.Warehouse{},
		&comp.RandomSprite{Rand: randSprite},
	)
	return e
}

func (f *EntityFactory) createPath(pos image.Point, t terr.Terrain, randSprite uint16) ecs.Entity {
	e := f.pathBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.Path{Haulers: []comp.HaulerEntry{}},
		&comp.RandomSprite{Rand: randSprite},
	)
	return e
}

func (f *EntityFactory) createProduction(pos image.Point, t terr.Terrain, prod *terr.Production, randSprite uint16) ecs.Entity {
	update := f.update.Get()
	e := f.productionBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.UpdateTick{Tick: rand.Int63n(update.Interval)},
		&comp.Production{Resource: prod.Resource, Amount: 0, Countdown: update.Countdown},
		&comp.RandomSprite{Rand: randSprite},
	)
	return e
}

func (f *EntityFactory) create(pos image.Point, t terr.Terrain, randSprite uint16) ecs.Entity {
	props := &terr.Properties[t]
	var e ecs.Entity
	if props.TerrainBits.Contains(terr.IsWarehouse) {
		e = f.createWarehouse(pos, t, randSprite)
	} else if props.TerrainBits.Contains(terr.IsPath) {
		e = f.createPath(pos, t, randSprite)
	} else {
		prod := props.Production
		if prod.MaxProduction == 0 {
			e = f.createLandUse(pos, t, randSprite)
		} else {
			e = f.createProduction(pos, t, &prod, randSprite)
		}
	}

	if props.BuildRadius > 0 {
		f.radiusMapper.Assign(e, &comp.BuildRadius{Radius: props.BuildRadius})
	}

	hasConsumption := false
	for _, c := range props.Consumption {
		if c > 0 {
			hasConsumption = true
			break
		}
	}
	if hasConsumption {
		cons := make([]uint8, len(props.Consumption))
		copy(cons, props.Consumption)
		f.consumptionMapper.Assign(e, &comp.Consumption{
			Amount:    cons,
			Countdown: make([]int16, len(props.Consumption)),
		})
	}
	if props.Population > 0 {
		f.populationMapper.Assign(e, &comp.Population{Pop: props.Population})
	}
	if props.PopulationSupport.MaxPopulation > 0 {
		f.populationSupportMapper.Add(e)
	}

	return e
}

// Set creates an entity of the given terrain type, placing it in the world and updating the game grids.
func (f *EntityFactory) Set(world *ecs.World, x, y int, value terr.Terrain, randSprite uint16) ecs.Entity {
	if !terr.Properties[value].TerrainBits.Contains(terr.IsTerrain) {
		f.landUse.Get().Set(x, y, value)
		e := f.create(image.Pt(x, y), value, randSprite)
		f.landUseEntities.Get().Set(x, y, e)

		rad := terr.Properties[value].BuildRadius
		if rad > 0 {
			f.SetBuildable(x, y, int(rad), true)
		}
		return e
	}
	t := f.terrain.Get()
	tE := f.terrainEntities.Get()

	eHere := tE.Get(x, y)
	if !eHere.IsZero() {
		world.RemoveEntity(eHere)
	}

	t.Set(x, y, value)
	e := f.create(image.Pt(x, y), value, randSprite)
	tE.Set(x, y, e)

	f.setNeighbor(t, tE, x-1, y)
	f.setNeighbor(t, tE, x+1, y)
	f.setNeighbor(t, tE, x, y-1)
	f.setNeighbor(t, tE, x, y+1)

	bounds := f.bounds.Get()
	bounds.AddPoint(image.Pt(x, y))

	return e
}

func (f *EntityFactory) setNeighbor(t *Terrain, tE *TerrainEntities, x, y int) {
	if t.Contains(x, y) && t.Get(x, y) == terr.Air {
		t.Set(x, y, terr.Buildable)
		e := f.create(image.Pt(x, y), terr.Buildable, 0)
		tE.Set(x, y, e)
	}
}

// RemoveLandUse removes land use from a given position, and updates the game grids.
func (f *EntityFactory) RemoveLandUse(world *ecs.World, x, y int) {
	landUse := f.landUse.Get()
	luHere := landUse.Get(x, y)
	if luHere == terr.Air {
		return
	}

	rad := terr.Properties[luHere].BuildRadius
	if rad > 0 {
		f.SetBuildable(x, y, int(rad), false)
	}

	luE := f.landUseEntities.Get()
	world.RemoveEntity(luE.Get(x, y))
	luE.Set(x, y, ecs.Entity{})
	landUse.Set(x, y, terr.Air)
}

// SetBuildable updates the build-ability grid.
// Only used for initialization, not required when using [EntityFactory.Set] or [EntityFactory.RemoveLandUse].
func (f *EntityFactory) SetBuildable(x, y, r int, build bool) {
	var add = 1
	if !build {
		add = -1
	}
	buildable := f.buildable.Get()
	r2 := r * r
	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			xx, yy := x+dx, y+dy
			if !buildable.Contains(xx, yy) || dx*dx+dy*dy > r2 {
				continue
			}
			v := buildable.Get(xx, yy)
			buildable.Set(xx, yy, uint16(int(v)+add))
		}
	}
}
