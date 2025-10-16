package sys

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Build system.
type Build struct {
	time            ecs.Resource[res.GameTick]
	rules           ecs.Resource[res.Rules]
	view            ecs.Resource[res.View]
	terrain         ecs.Resource[res.Terrain]
	terrainEntities ecs.Resource[res.TerrainEntities]
	landUse         ecs.Resource[res.LandUse]
	buildable       ecs.Resource[res.Buildable]
	stock           ecs.Resource[res.Stock]
	selection       ecs.Resource[res.Selection]
	update          ecs.Resource[res.UpdateInterval]
	ui              ecs.Resource[res.UI]
	factory         ecs.Resource[res.EntityFactory]
	editor          ecs.Resource[res.EditorMode]
	randTerrains    ecs.Resource[res.RandomTerrains]

	radiusFilter    *ecs.Filter2[comp.Tile, comp.BuildRadius]
	warehouseFilter *ecs.Filter1[comp.Warehouse]
}

// Initialize the system
func (s *Build) Initialize(world *ecs.World) {
	s.time = ecs.NewResource[res.GameTick](world)
	s.rules = ecs.NewResource[res.Rules](world)
	s.view = ecs.NewResource[res.View](world)
	s.terrain = ecs.NewResource[res.Terrain](world)
	s.terrainEntities = ecs.NewResource[res.TerrainEntities](world)
	s.landUse = ecs.NewResource[res.LandUse](world)
	s.buildable = ecs.NewResource[res.Buildable](world)
	s.stock = ecs.NewResource[res.Stock](world)
	s.selection = ecs.NewResource[res.Selection](world)
	s.update = ecs.NewResource[res.UpdateInterval](world)
	s.ui = ecs.NewResource[res.UI](world)
	s.factory = ecs.NewResource[res.EntityFactory](world)
	s.editor = ecs.NewResource[res.EditorMode](world)
	s.randTerrains = ecs.NewResource[res.RandomTerrains](world)

	s.radiusFilter = ecs.NewFilter2[comp.Tile, comp.BuildRadius](world)
	s.warehouseFilter = ecs.NewFilter1[comp.Warehouse](world)
}

// Update the system
func (s *Build) Update(world *ecs.World) {
	ui := s.ui.Get()
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		ui.ClearSelection()
		return
	}
	isEditor := s.editor.Get().IsEditor
	if s.checkAbort(isEditor) {
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
	randTerr := s.randTerrains.Get()
	landUse := s.landUse.Get()

	if !isEditor {
		if !stock.CanPay(p.BuildCost) {
			ui.SetStatusLabel("Not enough resources.")
			return
		}
		if p.TerrainBits.Contains(terr.CanBuild) && !p.TerrainBits.Contains(terr.CanBuy) {
			if randTerr.TotalPlaced >= randTerr.TotalAvailable {
				ui.SetStatusLabel("No more random terrains available.")
				return
			}
		}
	}

	if sel.BuildType == terr.Bulldoze {
		luHere := landUse.Get(cursor.X, cursor.Y)
		luProps := &terr.Properties[luHere]

		if luProps.TerrainBits.Contains(terr.IsWarehouse) && s.isLastWarehouse(stock, luHere) {
			ui.SetStatusLabel("Can't destroy last warehouse.")
			return
		}

		if luProps.TerrainBits.Contains(terr.CanBuild) {
			fac.RemoveLandUse(world, cursor.X, cursor.Y)

			if !isEditor {
				stock.Pay(p.BuildCost)
			}
			ui.ReplaceButton(stock, rules, randTerr, s.time.Get().RenderTick, image.Pt(x, y))
		}
		return
	}

	if !isEditor {
		if p.Population > 0 && stock.Population+int(p.Population) > stock.MaxPopulation {
			ui.SetStatusLabel("Population limit reached.")
			return
		}
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

	if !isEditor {
		stock.Pay(p.BuildCost)
	}
	ui.ReplaceButton(stock, rules, randTerr, s.time.Get().RenderTick, image.Pt(x, y))
}

// Finalize the system
func (s *Build) Finalize(world *ecs.World) {}

func (s *Build) checkAbort(isEditor bool) bool {
	if isEditor {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			return true
		}
	} else {
		if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			return true
		}
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
