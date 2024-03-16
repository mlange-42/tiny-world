package game

import (
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/menu"
	"github.com/mlange-42/tiny-world/game/render"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/res/achievements"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/save"
	"github.com/mlange-42/tiny-world/game/sys"
	"github.com/mlange-42/tiny-world/game/terr"
)

const TPS = 60
const saveFolder = "save"
const mapsFolder = "maps"

var GameData embed.FS

func Run(data embed.FS) {
	GameData = data

	resource.Prepare(GameData, "data/json/resources.json")
	terr.Prepare(GameData, "data/json/terrain.json")

	game := NewGame(nil)
	runMenu(&game, 0)

	game.Initialize()
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}

func run(g *Game, name string, mapLoc save.MapLocation, load save.LoadType, isEditor bool) {
	if err := runGame(g, load, name, mapLoc, "paper", isEditor); err != nil {
		panic(err)
	}
}

func runMenu(g *Game, tab int) {
	ebiten.SetVsyncEnabled(true)
	g.Model = model.New()

	ecs.AddResource(&g.Model.World, &g.Screen)

	sprites := res.NewSprites(GameData, "data/gfx", "paper")
	ecs.AddResource(&g.Model.World, &sprites)

	achievements := achievements.New(&g.Model.World, GameData, "data/json/achievements.json", "user/achievements.json")

	fonts := res.NewFonts(GameData)
	ui := menu.NewUI(GameData, saveFolder, mapsFolder, tab, &sprites, &fonts, achievements,
		func(name string, mapLoc save.MapLocation, load save.LoadType, isEditor bool) {
			run(g, name, mapLoc, load, isEditor)
		},
		func(tab int) {
			runMenu(g, tab)
		},
	)

	ecs.AddResource(&g.Model.World, &ui)

	g.Model.AddSystem(&menu.UpdateUI{})
	g.Model.AddUISystem(&menu.DrawUI{})

	g.Model.Initialize()
}

func runGame(g *Game, load save.LoadType, name string, mapLoc save.MapLocation, tileSet string, isEditor bool) error {
	ebiten.SetVsyncEnabled(true)

	g.Model = model.New()

	_ = ecs.ComponentID[comp.CardAnimation](&g.Model.World)

	// =========== Resources ===========

	rules := res.NewRules(GameData, "data/json/rules.json")
	ecs.AddResource(&g.Model.World, &rules)

	gameSpeed := res.GameSpeed{
		MinSpeed: -2,
		MaxSpeed: 3,
	}
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

	buildable := res.NewBuildable(rules.WorldSize, rules.WorldSize)
	ecs.AddResource(&g.Model.World, &buildable)

	selection := res.Selection{}
	ecs.AddResource(&g.Model.World, &selection)

	bounds := res.WorldBounds{}
	ecs.AddResource(&g.Model.World, &bounds)

	editor := res.EditorMode{IsEditor: isEditor}
	ecs.AddResource(&g.Model.World, &editor)

	randomTerrains := res.RandomTerrains{
		TotalAvailable: rules.InitialRandomTerrains,
	}
	ecs.AddResource(&g.Model.World, &randomTerrains)

	update := res.UpdateInterval{
		Interval:  TPS,
		Countdown: 60,
	}
	ecs.AddResource(&g.Model.World, &update)

	sprites := res.NewSprites(GameData, "data/gfx", tileSet)
	ecs.AddResource(&g.Model.World, &sprites)

	view := res.NewView(sprites.TileWidth, sprites.TileHeight)
	ecs.AddResource(&g.Model.World, &view)

	production := res.NewProduction()
	ecs.AddResource(&g.Model.World, &production)

	stock := res.NewStock(rules.InitialResources)
	ecs.AddResource(&g.Model.World, &stock)

	ecs.AddResource(&g.Model.World, &g.Screen)
	ecs.AddResource(&g.Model.World, &g.Mouse)

	saveEvent := res.SaveEvent{}
	ecs.AddResource(&g.Model.World, &saveEvent)

	fonts := res.NewFonts(GameData)
	ecs.AddResource(&g.Model.World, &fonts)

	factory := res.NewEntityFactory(&g.Model.World)
	ecs.AddResource(&g.Model.World, &factory)

	achievements := achievements.New(&g.Model.World, GameData, "data/json/achievements.json", "user/achievements.json")
	ecs.AddResource(&g.Model.World, achievements)

	// =========== Systems ===========

	if load == save.LoadTypeGame {
		g.Model.AddSystem(&sys.InitTerrainLoaded{})
	} else if load == save.LoadTypeMap {
		g.Model.AddSystem(&sys.InitTerrainMap{
			FS:        GameData,
			MapFolder: mapsFolder,
			Map:       mapLoc,
		})
	} else {
		g.Model.AddSystem(&sys.InitTerrain{})
	}
	g.Model.AddSystem(&sys.InitUI{})

	g.Model.AddSystem(&sys.Tick{})
	g.Model.AddSystem(&sys.UpdateProduction{})
	g.Model.AddSystem(&sys.UpdatePopulation{})
	g.Model.AddSystem(&sys.DoProduction{})
	g.Model.AddSystem(&sys.DoConsumption{})
	g.Model.AddSystem(&sys.Haul{})
	g.Model.AddSystem(&sys.UpdateStats{})
	g.Model.AddSystem(&sys.RemoveMarkers{
		MaxTime: TPS,
	})

	g.Model.AddSystem(&sys.Build{})
	g.Model.AddSystem(&sys.AssignHaulers{})
	g.Model.AddSystem(&sys.Achievements{
		PlayerFile: "user/achievements.json",
	})

	g.Model.AddSystem(&sys.PanAndZoom{
		PanButton:        ebiten.MouseButton1,
		ZoomInKey:        '+',
		ZoomOutKey:       '-',
		KeyboardPanSpeed: 4,
		MinZoom:          0.25,
		MaxZoom:          4,
	})

	g.Model.AddSystem(&sys.UpdateUI{})
	g.Model.AddSystem(&sys.Cheats{})
	g.Model.AddSystem(&sys.SaveGame{
		SaveFolder:   "save",
		MapFolder:    "maps",
		Name:         name,
		MainMenuFunc: func() { runMenu(g, 0) },
	})
	g.Model.AddSystem(&sys.GameControls{
		PauseKey:      ebiten.KeySpace,
		SlowerKey:     '[',
		FasterKey:     ']',
		FullscreenKey: ebiten.KeyF11,
	})

	// =========== UI Systems ===========

	g.Model.AddUISystem(&render.CenterView{})
	g.Model.AddUISystem(&render.Terrain{})
	g.Model.AddUISystem(&render.Markers{
		MinOffset: view.TileHeight * 2,
		MaxOffset: view.TileHeight*2 + 30,
		Duration:  TPS,
	})
	g.Model.AddUISystem(&render.UI{})
	g.Model.AddUISystem(&render.CardAnimation{
		MaxOffset: 200,
		Duration:  TPS / 4,
	})

	// =========== Load game ===========
	if load == save.LoadTypeGame {
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
