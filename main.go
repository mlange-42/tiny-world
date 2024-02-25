package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche-model/model"
	serde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/render"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/sys"
)

func main() {
	loadGame := len(os.Args) == 2

	ebiten.SetVsyncEnabled(true)
	g := game.NewGame(model.New())

	// =========== Resources ===========

	rules := res.Rules{
		AllowStroke:         false,
		AllowReplaceTerrain: false,
		AllowRemoveNatural:  false,
		AllowRemoveBuilt:    true,

		WorldSize:      128,
		RandomTerrains: 5,

		InitialResources:  [3]int{25, 25, 25},
		StockPerWarehouse: [3]int{25, 25, 25},
	}
	ecs.AddResource(&g.Model.World, &rules)

	terrain := res.NewTerrain(rules.WorldSize, rules.WorldSize)
	ecs.AddResource(&g.Model.World, &terrain)

	landUse := res.NewLandUse(rules.WorldSize, rules.WorldSize)
	ecs.AddResource(&g.Model.World, &landUse)

	landUseEntities := res.LandUseEntities{Grid: res.NewGrid[ecs.Entity](rules.WorldSize, rules.WorldSize)}
	ecs.AddResource(&g.Model.World, &landUseEntities)

	selection := res.Selection{}
	ecs.AddResource(&g.Model.World, &selection)

	update := res.UpdateInterval{
		Interval:  60,
		Countdown: 60,
	}
	ecs.AddResource(&g.Model.World, &update)

	view := res.NewView(48, 24)
	ecs.AddResource(&g.Model.World, &view)

	production := res.Production{}

	ecs.AddResource(&g.Model.World, &production)

	stock := res.Stock{Res: rules.InitialResources}
	ecs.AddResource(&g.Model.World, &stock)

	ecs.AddResource(&g.Model.World, &g.Screen)

	sprites := res.NewSprites("./assets/sprites")
	ecs.AddResource(&g.Model.World, &sprites)

	fonts := res.NewFonts()
	ecs.AddResource(&g.Model.World, &fonts)

	ui := res.NewUI(&selection, fonts.Default, &sprites, view.TileWidth)
	ecs.AddResource(&g.Model.World, &ui)

	factory := res.NewEntityFactory(&g.Model.World)
	ecs.AddResource(&g.Model.World, &factory)

	// =========== Systems ===========

	if !loadGame {
		g.Model.AddSystem(&sys.InitTerrain{})
	}

	g.Model.AddSystem(&sys.UpdateProduction{})
	g.Model.AddSystem(&sys.DoProduction{})
	g.Model.AddSystem(&sys.DoConsumption{})
	g.Model.AddSystem(&sys.UpdateStats{})
	g.Model.AddSystem(&sys.RemoveMarkers{
		MaxTime: 180,
	})

	g.Model.AddSystem(&sys.Build{})

	g.Model.AddSystem(&sys.PanAndZoom{
		PanButton: ebiten.MouseButton1,
	})

	g.Model.AddSystem(&sys.UpdateUI{})
	g.Model.AddSystem(&sys.SaveGame{
		Path: "./save/autosave.json",
	})

	// =========== UI Systems ===========

	g.Model.AddUISystem(&render.CenterView{})
	g.Model.AddUISystem(&render.Terrain{})
	//g.Model.AddUISystem(&render.Path{})
	g.Model.AddUISystem(&render.HaulerPaths{})
	g.Model.AddUISystem(&render.Markers{
		MinOffset: view.TileHeight * 2,
		MaxOffset: 250,
		Duration:  180,
	})
	g.Model.AddUISystem(&render.UI{})

	// =========== Load game ===========
	if loadGame {
		load(&g.Model.World, os.Args[1])
		selection.Reset()
	}

	// =========== Run ===========

	g.Initialize()
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}

func load(world *ecs.World, path string) {
	_ = ecs.ComponentID[comp.Tile](world)
	_ = ecs.ComponentID[comp.UpdateTick](world)
	_ = ecs.ComponentID[comp.Consumption](world)
	_ = ecs.ComponentID[comp.Production](world)
	_ = ecs.ComponentID[comp.Warehouse](world)
	_ = ecs.ComponentID[comp.Hauler](world)
	_ = ecs.ComponentID[comp.ProductionMarker](world)

	js, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = serde.Deserialize(js, world)
	if err != nil {
		panic(err)
	}
}
