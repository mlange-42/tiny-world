package res

import (
	"image"
	"math"
	"math/rand"

	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/terr"
)

// EntityFactory is a helper to create game entities.
type EntityFactory struct {
	landUseBuilder    *ecs.Map4[comp.Tile, comp.Terrain, comp.UpdateTick, comp.RandomSprite]
	productionBuilder *ecs.Map5[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Production, comp.RandomSprite]
	pathBuilder       *ecs.Map4[comp.Tile, comp.Terrain, comp.Path, comp.RandomSprite]

	radiusMapper            *ecs.Map1[comp.BuildRadius]
	consumptionMapper       *ecs.Map1[comp.Consumption]
	populationMapper        *ecs.Map1[comp.Population]
	populationSupportMapper *ecs.Map1[comp.PopulationSupport]
	unlockMapper            *ecs.Map1[comp.UnlocksTerrain]
	warehouseMapper         *ecs.Map1[comp.Warehouse]

	terrain         ecs.Resource[Terrain]
	terrainEntities ecs.Resource[TerrainEntities]
	landUse         ecs.Resource[LandUse]
	landUseEntities ecs.Resource[LandUseEntities]
	buildable       ecs.Resource[Buildable]
	bounds          ecs.Resource[WorldBounds]

	update ecs.Resource[UpdateInterval]
}

// NewEntityFactory creates a new EntityFactory for a given world.
func NewEntityFactory(world *ecs.World) EntityFactory {
	return EntityFactory{
		landUseBuilder:    ecs.NewMap4[comp.Tile, comp.Terrain, comp.UpdateTick, comp.RandomSprite](world),
		productionBuilder: ecs.NewMap5[comp.Tile, comp.Terrain, comp.UpdateTick, comp.Production, comp.RandomSprite](world),
		pathBuilder:       ecs.NewMap4[comp.Tile, comp.Terrain, comp.Path, comp.RandomSprite](world),

		radiusMapper:            ecs.NewMap1[comp.BuildRadius](world),
		consumptionMapper:       ecs.NewMap1[comp.Consumption](world),
		populationMapper:        ecs.NewMap1[comp.Population](world),
		populationSupportMapper: ecs.NewMap1[comp.PopulationSupport](world),
		warehouseMapper:         ecs.NewMap1[comp.Warehouse](world),
		unlockMapper:            ecs.NewMap1[comp.UnlocksTerrain](world),

		terrain:         ecs.NewResource[Terrain](world),
		terrainEntities: ecs.NewResource[TerrainEntities](world),
		landUse:         ecs.NewResource[LandUse](world),
		landUseEntities: ecs.NewResource[LandUseEntities](world),
		buildable:       ecs.NewResource[Buildable](world),
		bounds:          ecs.NewResource[WorldBounds](world),

		update: ecs.NewResource[UpdateInterval](world),
	}
}

func (f *EntityFactory) createLandUse(pos image.Point, t terr.Terrain, randSprite uint16) ecs.Entity {
	e := f.landUseBuilder.NewEntity(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.UpdateTick{Tick: rand.Int63n(f.update.Get().Interval)},
		&comp.RandomSprite{Rand: randSprite},
	)
	return e
}

func (f *EntityFactory) createPath(pos image.Point, t terr.Terrain, randSprite uint16) ecs.Entity {
	e := f.pathBuilder.NewEntity(
		&comp.Tile{Point: pos},
		&comp.Terrain{Terrain: t},
		&comp.Path{Haulers: []comp.HaulerEntry{}},
		&comp.RandomSprite{Rand: randSprite},
	)
	return e
}

func (f *EntityFactory) createProduction(pos image.Point, t terr.Terrain, prod *terr.Production, randSprite uint16) ecs.Entity {
	update := f.update.Get()
	e := f.productionBuilder.NewEntity(
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
	if props.TerrainBits.Contains(terr.IsPath) {
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
		f.radiusMapper.Add(e, &comp.BuildRadius{Radius: props.BuildRadius})
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
		f.consumptionMapper.Add(e, &comp.Consumption{
			Amount:    cons,
			Countdown: make([]int16, len(props.Consumption)),
		})
	}
	if props.Population > 0 {
		f.populationMapper.Add(e, &comp.Population{Pop: props.Population})
	}
	if props.PopulationSupport.MaxPopulation > 0 {
		f.populationSupportMapper.AddFn(e, nil)
	}
	if props.UnlocksTerrains > 0 {
		f.unlockMapper.AddFn(e, nil)
	}
	if props.TerrainBits.Contains(terr.IsWarehouse) {
		f.warehouseMapper.AddFn(e, nil)
	}

	return e
}

// Set creates an entity of the given terrain type, placing it in the world and updating the game grids.
func (f *EntityFactory) Set(world *ecs.World, x, y int, value terr.Terrain, randSprite uint16, randomize bool) ecs.Entity {
	if randomize {
		randSprite = uint16(rand.Int31n(math.MaxUint16))
	}
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

	if value != terr.Air && value != terr.Buildable {
		bounds := f.bounds.Get()
		bounds.AddPoint(image.Pt(x, y))
	}

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
