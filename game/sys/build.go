package sys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Build system.
type Build struct {
	rules           generic.Resource[res.Rules]
	view            generic.Resource[res.View]
	terrain         generic.Resource[res.Terrain]
	landUse         generic.Resource[res.LandUse]
	landUseEntities generic.Resource[res.LandUseEntities]
	stock           generic.Resource[res.Stock]
	selection       generic.Resource[res.Selection]
	update          generic.Resource[res.UpdateInterval]
	ui              generic.Resource[res.UI]
	factory         generic.Resource[res.EntityFactory]
}

// Initialize the system
func (s *Build) Initialize(world *ecs.World) {
	s.rules = generic.NewResource[res.Rules](world)
	s.view = generic.NewResource[res.View](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.landUseEntities = generic.NewResource[res.LandUseEntities](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.selection = generic.NewResource[res.Selection](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.ui = generic.NewResource[res.UI](world)
	s.factory = generic.NewResource[res.EntityFactory](world)
}

// Update the system
func (s *Build) Update(world *ecs.World) {
	rules := s.rules.Get()
	sel := s.selection.Get()
	ui := s.ui.Get()
	fac := s.factory.Get()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sel.Reset()
	}

	x, y := ebiten.CursorPosition()
	if ui.MouseInside(x, y) {
		return
	}

	mouseFn := inpututil.IsMouseButtonJustPressed
	if rules.AllowStroke {
		mouseFn = ebiten.IsMouseButtonPressed
	}

	p := &terr.Properties[sel.BuildType]
	if !p.CanBuild ||
		!(mouseFn(ebiten.MouseButton0) ||
			mouseFn(ebiten.MouseButton2)) {
		return
	}

	stock := s.stock.Get()
	view := s.view.Get()
	landUse := s.landUse.Get()
	landUseE := s.landUseEntities.Get()
	mx, my := view.ScreenToGlobal(x, y)
	cursor := view.GlobalToTile(mx, my)

	remove := mouseFn(ebiten.MouseButton2)
	if remove {
		if p.IsTerrain {
			sel.Reset()
			return
		}
		luHere := landUse.Get(cursor.X, cursor.Y)
		canBuy := p.CanBuy
		if luHere == sel.BuildType &&
			((rules.AllowRemoveBuilt && canBuy) || (rules.AllowRemoveNatural && !canBuy) &&
				stock.CanPay(p.BuildCost)) {

			world.RemoveEntity(landUseE.Get(cursor.X, cursor.Y))
			landUseE.Set(cursor.X, cursor.Y, ecs.Entity{})
			landUse.Set(cursor.X, cursor.Y, terr.Air)

			stock.Pay(p.BuildCost)
			ui.ReplaceButton(stock)
		} else {
			sel.Reset()
		}
		return
	}
	if landUse.Get(cursor.X, cursor.Y) != terr.Air {
		return
	}

	if !stock.CanPay(p.BuildCost) {
		return
	}

	terrain := s.terrain.Get()
	terrHere := terrain.Get(cursor.X, cursor.Y)
	if p.IsTerrain {
		if rules.AllowReplaceTerrain {
			if !p.BuildOnFree.Contains(terrHere) {
				return
			}
		} else {
			if !p.BuildOn.Contains(terrHere) {
				return
			}
		}
		terrain.Set(cursor.X, cursor.Y, sel.BuildType)
	} else {
		if !p.BuildOn.Contains(terrHere) {
			return
		}
		e := fac.Create(cursor, sel.BuildType)
		landUseE.Set(cursor.X, cursor.Y, e)
		landUse.Set(cursor.X, cursor.Y, sel.BuildType)
	}

	stock.Pay(p.BuildCost)
	ui.ReplaceButton(stock)
}

// Finalize the system
func (s *Build) Finalize(world *ecs.World) {}
