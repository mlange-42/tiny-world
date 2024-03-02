package game

import (
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/menu"
	"github.com/mlange-42/tiny-world/game/render"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/save"
	"github.com/mlange-42/tiny-world/game/sys"
	"github.com/mlange-42/tiny-world/game/terr"
)

const saveFolder = "save"

var gameData embed.FS

func Run(data embed.FS) {
	gameData = data

	runMenu()
}

func runMenu() {
	ebiten.SetVsyncEnabled(true)
	g := NewGame(model.New())

	ecs.AddResource(&g.Model.World, &g.Screen)

	fonts := res.NewFonts(gameData)
	ui := menu.NewUI(saveFolder, fonts.Default, func(name string, load bool) {
		run(&g, name, load)
	})

	ecs.AddResource(&g.Model.World, &ui)

	g.Model.AddSystem(&menu.UpdateUI{})
	g.Model.AddUISystem(&menu.DrawUI{})

	g.Initialize()
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}

func runGame(g *Game, loadGame bool, name, tileSet string) error {
	ebiten.SetVsyncEnabled(true)

	g.Model = model.New()

	resource.Prepare(gameData, "data/json/resources.json")
	terr.Prepare(gameData, "data/json/terrain.json")

	// =========== Resources ===========

	rules := res.NewRules(gameData, "data/json/rules.json")
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

	sprites := res.NewSprites(gameData, "data/gfx", tileSet)
	ecs.AddResource(&g.Model.World, &sprites)

	view := res.NewView(sprites.TileWidth, sprites.TileHeight)
	ecs.AddResource(&g.Model.World, &view)

	production := res.NewProduction()
	ecs.AddResource(&g.Model.World, &production)

	stock := res.NewStock(rules.InitialResources)
	ecs.AddResource(&g.Model.World, &stock)

	ecs.AddResource(&g.Model.World, &g.Screen)

	fonts := res.NewFonts(gameData)
	ecs.AddResource(&g.Model.World, &fonts)

	ui := res.NewUI(&selection, fonts.Default, &sprites)
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
		Folder: "save",
		Name:   name,
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
		err := save.LoadWorld(&g.Model.World, saveFolder, name)
		if err != nil {
			return err
		}
		selection.Reset()

		view.TileWidth = sprites.TileWidth
		view.TileHeight = sprites.TileHeight
	}

	// =========== Run ===========

	g.Model.Initialize()

	return nil
}
