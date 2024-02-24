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
	"github.com/mlange-42/tiny-world/game/terr"
)

const (
	WORLD_SIZE = 128

	FOOD   = 25
	WOOD   = 10
	STONES = 5

	RANDOM_TERRAINS = 5
)

func main() {
	loadGame := len(os.Args) == 2

	ebiten.SetVsyncEnabled(true)
	g := game.NewGame(model.New())

	// =========== Resources ===========

	terrain := res.Terrain{Grid: res.NewGrid[terr.Terrain](WORLD_SIZE, WORLD_SIZE)}
	ecs.AddResource(&g.Model.World, &terrain)

	landUse := res.LandUse{Grid: res.NewGrid[terr.Terrain](WORLD_SIZE, WORLD_SIZE)}
	ecs.AddResource(&g.Model.World, &landUse)

	landUseEntities := res.LandUseEntities{Grid: res.NewGrid[ecs.Entity](WORLD_SIZE, WORLD_SIZE)}
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
	stock := res.Stock{Res: [3]int{FOOD, WOOD, STONES}}
	ecs.AddResource(&g.Model.World, &production)
	ecs.AddResource(&g.Model.World, &stock)

	ecs.AddResource(&g.Model.World, &g.Screen)

	sprites := res.NewSprites("./assets/sprites")
	ecs.AddResource(&g.Model.World, &sprites)

	fonts := res.NewFonts()
	ecs.AddResource(&g.Model.World, &fonts)

	ui := res.NewUI(&selection, fonts.Default, &sprites, view.TileWidth, RANDOM_TERRAINS)
	ecs.AddResource(&g.Model.World, &ui)

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

	g.Model.AddSystem(&sys.Build{
		AllowStroke: true,
	})

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
	g.Model.AddUISystem(&render.Markers{
		MinOffset: view.TileHeight * 2,
		MaxOffset: 250,
		Duration:  180,
	})
	g.Model.AddUISystem(&render.UI{})

	// =========== Load game ===========
	if loadGame {
		load(&g.Model.World, os.Args[1])
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
