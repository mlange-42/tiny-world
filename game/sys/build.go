package sys

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Build system.
type Build struct {
	AllowStroke bool

	view            generic.Resource[res.View]
	terrain         generic.Resource[res.Terrain]
	landUse         generic.Resource[res.LandUse]
	landUseEntities generic.Resource[res.LandUseEntities]
	selection       generic.Resource[res.Selection]
	update          generic.Resource[res.UpdateInterval]

	builder           generic.Map2[comp.Tile, comp.UpdateTick]
	productionBuilder generic.Map3[comp.Tile, comp.UpdateTick, comp.Production]
}

// Initialize the system
func (s *Build) Initialize(world *ecs.World) {
	s.view = generic.NewResource[res.View](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.landUseEntities = generic.NewResource[res.LandUseEntities](world)
	s.selection = generic.NewResource[res.Selection](world)
	s.update = generic.NewResource[res.UpdateInterval](world)

	s.builder = generic.NewMap2[comp.Tile, comp.UpdateTick](world)
	s.productionBuilder = generic.NewMap3[comp.Tile, comp.UpdateTick, comp.Production](world)
}

// Update the system
func (s *Build) Update(world *ecs.World) {
	sel := s.selection.Get()

	for i := range terr.Properties {
		p := &terr.Properties[i]
		if p.CanBuild && inpututil.IsKeyJustPressed(p.ShortKey) {
			sel.Build = terr.Terrain(i)
			fmt.Println("paint", p.Name)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sel.Build = terr.Air
		fmt.Println("paint nothing")
	}

	mouseFn := inpututil.IsMouseButtonJustPressed
	if s.AllowStroke {
		mouseFn = ebiten.IsMouseButtonPressed
	}

	p := &terr.Properties[sel.Build]
	if !p.CanBuild ||
		!(mouseFn(ebiten.MouseButton0) ||
			mouseFn(ebiten.MouseButton2)) {
		return
	}

	view := s.view.Get()
	landUse := s.landUse.Get()
	landUseE := s.landUseEntities.Get()
	mx, my := view.ScreenToGlobal(ebiten.CursorPosition())
	cursor := view.GlobalToTile(mx, my)

	remove := mouseFn(ebiten.MouseButton2)
	if remove {
		if p.IsTerrain {
			return
		}
		luHere := landUse.Get(cursor.X, cursor.Y)
		if luHere == sel.Build {
			world.RemoveEntity(landUseE.Get(cursor.X, cursor.Y))
			landUseE.Set(cursor.X, cursor.Y, ecs.Entity{})
			landUse.Set(cursor.X, cursor.Y, terr.Air)
		}
		return
	}

	terrain := s.terrain.Get()

	terrHere := terrain.Get(cursor.X, cursor.Y)
	if !p.BuildOn.Contains(terrHere) {
		return
	}
	if p.IsTerrain {
		terrain.Set(cursor.X, cursor.Y, sel.Build)
	} else {
		update := s.update.Get()
		if landUse.Get(cursor.X, cursor.Y) != terr.Air {
			return
		}
		prod := terr.Properties[sel.Build].Production.Produces
		var e ecs.Entity
		if prod == resource.EndResources {
			e = s.builder.NewWith(
				&comp.Tile{Point: cursor},
				&comp.UpdateTick{Tick: rand.Int63n(update.Interval)},
			)
		} else {
			e = s.productionBuilder.NewWith(
				&comp.Tile{Point: cursor},
				&comp.UpdateTick{Tick: rand.Int63n(update.Interval)},
				&comp.Production{Type: prod, Amount: 0},
			)
		}
		landUseE.Set(cursor.X, cursor.Y, e)
		landUse.Set(cursor.X, cursor.Y, sel.Build)
	}
}

// Finalize the system
func (s *Build) Finalize(world *ecs.World) {}
