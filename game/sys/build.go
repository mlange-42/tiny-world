package sys

import (
	"fmt"
	"image"

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
	time            generic.Resource[res.GameTick]
	rules           generic.Resource[res.Rules]
	view            generic.Resource[res.View]
	terrain         generic.Resource[res.Terrain]
	terrainEntities generic.Resource[res.TerrainEntities]
	landUse         generic.Resource[res.LandUse]
	buildable       generic.Resource[res.Buildable]
	stock           generic.Resource[res.Stock]
	selection       generic.Resource[res.Selection]
	update          generic.Resource[res.UpdateInterval]
	ui              generic.Resource[res.UI]
	factory         generic.Resource[res.EntityFactory]

	radiusFilter    generic.Filter2[comp.Tile, comp.BuildRadius]
	warehouseFilter generic.Filter1[comp.Warehouse]
}

// Initialize the system
func (s *Build) Initialize(world *ecs.World) {
	s.time = generic.NewResource[res.GameTick](world)
	s.rules = generic.NewResource[res.Rules](world)
	s.view = generic.NewResource[res.View](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.terrainEntities = generic.NewResource[res.TerrainEntities](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.buildable = generic.NewResource[res.Buildable](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.selection = generic.NewResource[res.Selection](world)
	s.update = generic.NewResource[res.UpdateInterval](world)
	s.ui = generic.NewResource[res.UI](world)
	s.factory = generic.NewResource[res.EntityFactory](world)

	s.radiusFilter = *generic.NewFilter2[comp.Tile, comp.BuildRadius]()
	s.warehouseFilter = *generic.NewFilter1[comp.Warehouse]()
}

// Update the system
func (s *Build) Update(world *ecs.World) {
	ui := s.ui.Get()
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		ui.ClearSelection()
		return
	}
	if s.checkAbort() {
		return
	}
	buildable := s.buildable.Get()
	sel := s.selection.Get()
	view := s.view.Get()
	x, y := ebiten.CursorPosition()
	mx, my := view.ScreenToGlobal(x, y)
	cursor := view.GlobalToTile(mx, my)

	p := &terr.Properties[sel.BuildType]
	if p.TerrainBits.Contains(terr.RequiresRange) && buildable.Get(cursor.X, cursor.Y) == 0 {
		ui.SetStatusLabel("Outside of controlled area.")
		return
	}

	fac := s.factory.Get()
	rules := s.rules.Get()
	stock := s.stock.Get()
	landUse := s.landUse.Get()

	if sel.BuildType == terr.Bulldoze {
		luHere := landUse.Get(cursor.X, cursor.Y)
		luProps := &terr.Properties[luHere]

		if luProps.TerrainBits.Contains(terr.IsWarehouse) && s.isLastWarehouse(stock, luHere) {
			ui.SetStatusLabel("Can't destroy last warehouse.")
			return
		}

		if luProps.TerrainBits.Contains(terr.CanBuild) {
			fac.RemoveLandUse(world, cursor.X, cursor.Y)

			stock.Pay(p.BuildCost)
			ui.ReplaceButton(stock, rules, s.time.Get().RenderTick, image.Pt(x, y))
		}
		return
	}

	if !stock.CanPay(p.BuildCost) {
		ui.SetStatusLabel("Not enough resources.")
		return
	}
	if p.Population > 0 && stock.Population+int(p.Population) > stock.MaxPopulation {
		ui.SetStatusLabel("Population limit reached.")
		return
	}

	terrain := s.terrain.Get()
	terrHere := terrain.Get(cursor.X, cursor.Y)
	luHere := landUse.Get(cursor.X, cursor.Y)
	if p.TerrainBits.Contains(terr.IsTerrain) {
		canBuild := luHere == terr.Air
		if terrHere == terr.Air {
			ui.SetStatusLabel("Can only add next to existing terrain.")
			return
		}
		canBuild = canBuild &&
			(p.BuildOn.Contains(terrHere) || (sel.AllowRemove && terrHere != sel.BuildType))
		if !canBuild {
			ui.SetStatusLabel("Terrain already occupied.")
			return
		}
		fac.Set(world, cursor.X, cursor.Y, sel.BuildType, sel.RandSprite, sel.Randomize)
	} else {
		if terrHere == terr.Air || terrHere == terr.Buildable {
			ui.SetStatusLabel("No terrain here.")
			return
		}
		if !p.BuildOn.Contains(terrHere) {
			ui.SetStatusLabel(fmt.Sprintf("Can't build this on %s", terr.Properties[terrHere].Name))
			return
		}

		luNatural := !terr.Properties[luHere].TerrainBits.Contains(terr.CanBuy)
		if luHere == terr.Air || (luNatural && p.TerrainBits.Contains(terr.CanBuy)) {
			if luHere != terr.Air {
				fac.RemoveLandUse(world, cursor.X, cursor.Y)
			}
			fac.Set(world, cursor.X, cursor.Y, sel.BuildType, sel.RandSprite, sel.Randomize)
		} else {
			ui.SetStatusLabel("Terrain already occupied.")
			return
		}
	}

	stock.Pay(p.BuildCost)
	ui.ReplaceButton(stock, rules, s.time.Get().RenderTick, image.Pt(x, y))
}

// Finalize the system
func (s *Build) Finalize(world *ecs.World) {}

func (s *Build) checkAbort() bool {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		return true
	}

	ui := s.ui.Get()
	x, y := ebiten.CursorPosition()
	if ui.MouseInside(x, y) {
		return true
	}

	sel := s.selection.Get()
	p := &terr.Properties[sel.BuildType]
	if sel.BuildType != terr.Bulldoze && !p.TerrainBits.Contains(terr.CanBuild) {
		return true
	}
	return false
}

func (s *Build) isLastWarehouse(stock *res.Stock, building terr.Terrain) bool {
	storage := terr.Properties[building].Storage
	for i := range resource.Properties {
		if stock.Cap[i] <= int(storage[i]) {
			return true
		}
	}
	return false
}
