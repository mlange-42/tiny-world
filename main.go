package main

import (
	"fmt"
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
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/sys"
	"github.com/mlange-42/tiny-world/game/terr"
	"github.com/spf13/cobra"
)

func main() {
	if err := command().Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(saveGame, tileSet string) {
	loadGame := saveGame != ""

	ebiten.SetVsyncEnabled(true)
	g := game.NewGame(model.New())

	resource.Prepare(data, "data/json/resources.json")
	terr.Prepare(data, "data/json/terrain.json")

	// =========== Resources ===========

	rules := res.NewRules(data, "data/json/rules.json")
	ecs.AddResource(&g.Model.World, &rules)

	gameSpeed := res.GameSpeed{}
	ecs.AddResource(&g.Model.World, &gameSpeed)

	gameTick := res.GameTick{}
	ecs.AddResource(&g.Model.World, &gameTick)

	terrain := res.NewTerrain(rules.WorldSize, rules.WorldSize)
	ecs.AddResource(&g.Model.World, &terrain)

	terrainEntities := res.TerrainEntities{Grid: res.NewGrid[ecs.Entity](rules.WorldSize, rules.WorldSize)}
	ecs.AddResource(&g.Model.World, &terrainEntities)

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

	sprites := res.NewSprites(data, "data/gfx", tileSet)
	ecs.AddResource(&g.Model.World, &sprites)

	view := res.NewView(sprites.TileWidth, sprites.TileHeight)
	ecs.AddResource(&g.Model.World, &view)

	production := res.NewProduction()
	ecs.AddResource(&g.Model.World, &production)

	stock := res.NewStock(rules.InitialResources)
	ecs.AddResource(&g.Model.World, &stock)

	ecs.AddResource(&g.Model.World, &g.Screen)

	fonts := res.NewFonts(data)
	ecs.AddResource(&g.Model.World, &fonts)

	ui := res.NewUI(&selection, fonts.Default, &sprites, sprites.TileWidth)
	ecs.AddResource(&g.Model.World, &ui)

	factory := res.NewEntityFactory(&g.Model.World)
	ecs.AddResource(&g.Model.World, &factory)

	// =========== Systems ===========

	if loadGame {
		g.Model.AddSystem(&sys.InitTerrainLoaded{})
	} else {
		g.Model.AddSystem(&sys.InitTerrain{})
	}

	g.Model.AddSystem(&sys.Tick{})
	g.Model.AddSystem(&sys.UpdateProduction{})
	g.Model.AddSystem(&sys.DoProduction{})
	g.Model.AddSystem(&sys.DoConsumption{})
	g.Model.AddSystem(&sys.Haul{})
	g.Model.AddSystem(&sys.UpdateStats{})
	g.Model.AddSystem(&sys.RemoveMarkers{
		MaxTime: 60,
	})

	g.Model.AddSystem(&sys.Build{})
	g.Model.AddSystem(&sys.AssignHaulers{})

	g.Model.AddSystem(&sys.PanAndZoom{
		PanButton: ebiten.MouseButton1,
	})

	g.Model.AddSystem(&sys.UpdateUI{})
	g.Model.AddSystem(&sys.Cheats{})
	g.Model.AddSystem(&sys.SaveGame{
		Path: "./save/autosave.json",
	})
	g.Model.AddSystem(&sys.Pause{
		PauseKey: ebiten.KeySpace,
	})

	// =========== UI Systems ===========

	g.Model.AddUISystem(&render.CenterView{})
	g.Model.AddUISystem(&render.Terrain{})
	//g.Model.AddUISystem(&render.HaulerPaths{})
	g.Model.AddUISystem(&render.Markers{
		MinOffset: view.TileHeight * 2,
		MaxOffset: view.TileHeight*2 + 30,
		Duration:  60,
	})
	g.Model.AddUISystem(&render.UI{})

	// =========== Load game ===========
	if loadGame {
		load(&g.Model.World, saveGame)
		selection.Reset()

		view.TileWidth = sprites.TileWidth
		view.TileHeight = sprites.TileHeight
	}

	// =========== Run ===========

	g.Initialize()
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}

func load(world *ecs.World, path string) {
	_ = ecs.ComponentID[comp.Tile](world)
	_ = ecs.ComponentID[comp.Terrain](world)
	_ = ecs.ComponentID[comp.UpdateTick](world)
	_ = ecs.ComponentID[comp.Consumption](world)
	_ = ecs.ComponentID[comp.Production](world)
	_ = ecs.ComponentID[comp.Warehouse](world)
	_ = ecs.ComponentID[comp.Path](world)
	_ = ecs.ComponentID[comp.Hauler](world)
	_ = ecs.ComponentID[comp.HaulerSprite](world)
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

func command() *cobra.Command {
	var tileSet, saveFile string
	root := &cobra.Command{
		Use:           "tiny-world",
		Short:         "A tiny, slow-paced world and colony building game.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(saveFile, tileSet)
		},
	}
	root.Flags().StringVarP(&tileSet, "tileset", "t", "paper", "Tileset to use.")
	root.Flags().StringVarP(&saveFile, "savefile", "s", "", "Savefile to load.")

	return root
}
