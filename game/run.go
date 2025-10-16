package game

import (
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark/ecs"
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
	g.App = app.New()

	ecs.AddResource(&g.App.World, &g.Screen)

	sprites := res.NewSprites(GameData, "data/gfx", "paper")
	ecs.AddResource(&g.App.World, &sprites)

	achievements := achievements.New(&g.App.World, GameData, "data/json/achievements.json", "user/achievements.json")

	fonts := res.NewFonts(GameData)
	ui := menu.NewUI(GameData, saveFolder, mapsFolder, tab, &sprites, &fonts, achievements,
		func(name string, mapLoc save.MapLocation, load save.LoadType, isEditor bool) {
			run(g, name, mapLoc, load, isEditor)
		},
		func(tab int) {
			runMenu(g, tab)
		},
	)

	ecs.AddResource(&g.App.World, &ui)

	g.App.AddSystem(&menu.UpdateUI{})
	g.App.AddUISystem(&menu.DrawUI{})

	g.App.Initialize()
}

func runGame(g *Game, load save.LoadType, name string, mapLoc save.MapLocation, tileSet string, isEditor bool) error {
	ebiten.SetVsyncEnabled(true)

	g.App = app.New()

	// Register components for deserialization,
	// where it does not happen in systems already.
	_ = ecs.ComponentID[comp.CardAnimation](&g.App.World)

	// =========== Resources ===========

	rules := res.NewRules(GameData, "data/json/rules.json")
	ecs.AddResource(&g.App.World, &rules)

	gameSpeed := res.GameSpeed{
		MinSpeed: -2,
		MaxSpeed: 3,
	}
	ecs.AddResource(&g.App.World, &gameSpeed)

	gameTick := res.GameTick{}
	ecs.AddResource(&g.App.World, &gameTick)

	terrain := res.NewTerrain(rules.WorldSize, rules.WorldSize)
	ecs.AddResource(&g.App.World, &terrain)

	terrainEntities := res.TerrainEntities{Grid: res.NewGrid[ecs.Entity](rules.WorldSize, rules.WorldSize)}
	ecs.AddResource(&g.App.World, &terrainEntities)

	landUse := res.NewLandUse(rules.WorldSize, rules.WorldSize)
	ecs.AddResource(&g.App.World, &landUse)

	landUseEntities := res.LandUseEntities{Grid: res.NewGrid[ecs.Entity](rules.WorldSize, rules.WorldSize)}
	ecs.AddResource(&g.App.World, &landUseEntities)

	buildable := res.NewBuildable(rules.WorldSize, rules.WorldSize)
	ecs.AddResource(&g.App.World, &buildable)

	selection := res.Selection{}
	ecs.AddResource(&g.App.World, &selection)

	bounds := res.WorldBounds{}
	ecs.AddResource(&g.App.World, &bounds)

	editor := res.EditorMode{IsEditor: isEditor}
	ecs.AddResource(&g.App.World, &editor)

	saveTime := res.SaveTime{}
	ecs.AddResource(&g.App.World, &saveTime)

	randomTerrains := res.RandomTerrains{
		TotalAvailable: rules.InitialRandomTerrains,
	}
	ecs.AddResource(&g.App.World, &randomTerrains)

	update := res.UpdateInterval{
		Interval:  TPS,
		Countdown: 60,
	}
	ecs.AddResource(&g.App.World, &update)

	sprites := res.NewSprites(GameData, "data/gfx", tileSet)
	ecs.AddResource(&g.App.World, &sprites)

	view := res.NewView(sprites.TileWidth, sprites.TileHeight)
	ecs.AddResource(&g.App.World, &view)

	production := res.NewProduction()
	ecs.AddResource(&g.App.World, &production)

	stock := res.NewStock(rules.InitialResources)
	ecs.AddResource(&g.App.World, &stock)

	ecs.AddResource(&g.App.World, &g.Screen)
	ecs.AddResource(&g.App.World, &g.Mouse)

	saveEvent := res.SaveEvent{}
	ecs.AddResource(&g.App.World, &saveEvent)

	fonts := res.NewFonts(GameData)
	ecs.AddResource(&g.App.World, &fonts)

	factory := res.NewEntityFactory(&g.App.World)
	ecs.AddResource(&g.App.World, &factory)

	achievements := achievements.New(&g.App.World, GameData, "data/json/achievements.json", "user/achievements.json")
	ecs.AddResource(&g.App.World, achievements)

	// =========== Systems ===========

	if load == save.LoadTypeGame {
		g.App.AddSystem(&sys.InitTerrainLoaded{})
	} else if load == save.LoadTypeMap {
		g.App.AddSystem(&sys.InitTerrainMap{
			FS:        GameData,
			MapFolder: mapsFolder,
			Map:       mapLoc,
		})
	} else {
		g.App.AddSystem(&sys.InitTerrain{})
	}
	g.App.AddSystem(&sys.InitUI{})

	g.App.AddSystem(&sys.Tick{})
	g.App.AddSystem(&sys.UpdateProduction{})
	g.App.AddSystem(&sys.UpdatePopulation{})
	g.App.AddSystem(&sys.DoProduction{})
	g.App.AddSystem(&sys.DoConsumption{})
	g.App.AddSystem(&sys.Haul{})
	g.App.AddSystem(&sys.UpdateStats{})
	g.App.AddSystem(&sys.RemoveMarkers{
		MaxTime: TPS,
	})

	g.App.AddSystem(&sys.Build{})
	g.App.AddSystem(&sys.AssignHaulers{})
	g.App.AddSystem(&sys.Achievements{
		PlayerFile: "user/achievements.json",
	})

	g.App.AddSystem(&sys.PanAndZoom{
		PanButton:        ebiten.MouseButton1,
		ZoomInKey:        '+',
		ZoomOutKey:       '-',
		KeyboardPanSpeed: 4,
		MinZoom:          0.25,
		MaxZoom:          4,
	})

	g.App.AddSystem(&sys.UpdateUI{})
	g.App.AddSystem(&sys.Cheats{})
	g.App.AddSystem(&sys.SaveGame{
		SaveFolder:   "save",
		MapFolder:    "maps",
		Name:         name,
		MainMenuFunc: func() { runMenu(g, 0) },
	})
	g.App.AddSystem(&sys.GameControls{
		PauseKey:      ebiten.KeySpace,
		SlowerKey:     '[',
		FasterKey:     ']',
		FullscreenKey: ebiten.KeyF11,
	})

	// =========== UI Systems ===========

	g.App.AddUISystem(&render.CenterView{})
	g.App.AddUISystem(&render.Terrain{})
	g.App.AddUISystem(&render.Markers{
		MinOffset: view.TileHeight * 2,
		MaxOffset: view.TileHeight*2 + 30,
		Duration:  TPS,
	})
	g.App.AddUISystem(&render.UI{})
	g.App.AddUISystem(&render.CardAnimation{
		MaxOffset: 200,
		Duration:  TPS / 4,
	})

	// =========== Load game ===========
	if load == save.LoadTypeGame {
		err := save.LoadWorld(&g.App.World, saveFolder, name)
		if err != nil {
			return err
		}
		selection.Reset()

		view.TileWidth = sprites.TileWidth
		view.TileHeight = sprites.TileHeight
	}

	addRepl(g.App)

	// =========== Run ===========

	g.App.Initialize()

	return nil
}
