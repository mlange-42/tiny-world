package res

import (
	"image"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

type EntityFactory struct {
	landUseBuilder    generic.Map2[comp.Tile, comp.UpdateTick]
	productionBuilder generic.Map4[comp.Tile, comp.UpdateTick, comp.Production, comp.Consumption]

	update generic.Resource[UpdateInterval]
}

func NewEntityFactory(world *ecs.World) EntityFactory {
	return EntityFactory{
		landUseBuilder:    generic.NewMap2[comp.Tile, comp.UpdateTick](world),
		productionBuilder: generic.NewMap4[comp.Tile, comp.UpdateTick, comp.Production, comp.Consumption](world),
		update:            generic.NewResource[UpdateInterval](world),
	}
}

func (f *EntityFactory) createLandUse(pos image.Point) ecs.Entity {
	e := f.landUseBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.UpdateTick{Tick: rand.Int63n(f.update.Get().Interval)},
	)
	return e
}

func (f *EntityFactory) createProduction(pos image.Point, prod *terr.Production) ecs.Entity {
	update := f.update.Get()
	e := f.productionBuilder.NewWith(
		&comp.Tile{Point: pos},
		&comp.UpdateTick{Tick: rand.Int63n(update.Interval)},
		&comp.Production{Type: prod.Produces, Amount: 0, Countdown: update.Countdown},
		&comp.Consumption{Amount: prod.ConsumesFood, Countdown: update.Countdown},
	)
	return e
}

func (f *EntityFactory) Create(pos image.Point, t terr.Terrain) ecs.Entity {
	prod := terr.Properties[t].Production
	var e ecs.Entity
	if prod.Produces == resource.EndResources {
		e = f.createLandUse(pos)
	} else {
		e = f.createProduction(pos, &prod)
	}
	return e
}
