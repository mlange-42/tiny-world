package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game"
	"github.com/mlange-42/tiny-world/game/render"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/sys"
	"github.com/mlange-42/tiny-world/game/terr"
)

const (
	worldSize = 128
)

func main() {
	ebiten.SetVsyncEnabled(true)
	g := game.NewGame(model.New())

	// =========== Resources ===========

	ecs.AddResource(&g.Model.World, &g.Screen)

	terrain := res.Terrain{Grid: res.NewGrid[terr.Terrain](worldSize, worldSize)}
	ecs.AddResource(&g.Model.World, &terrain)

	landUse := res.LandUse{Grid: res.NewGrid[terr.Terrain](worldSize, worldSize)}
	ecs.AddResource(&g.Model.World, &landUse)

	landUseEntities := res.LandUseEntities{Grid: res.NewGrid[ecs.Entity](worldSize, worldSize)}
	ecs.AddResource(&g.Model.World, &landUseEntities)

	selection := res.Selection{}
	ecs.AddResource(&g.Model.World, &selection)

	update := res.UpdateInterval{
		Interval:  60,
		Countdown: 60,
	}
	ecs.AddResource(&g.Model.World, &update)

	sprites := res.NewSprites("./assets/sprites")
	ecs.AddResource(&g.Model.World, &sprites)

	fonts := res.NewFonts()
	ecs.AddResource(&g.Model.World, &fonts)

	hud := res.NewHUD(fonts.Default)
	ui := res.NewUI(&selection, fonts.Default, &sprites)
	ecs.AddResource(&g.Model.World, &hud)
	ecs.AddResource(&g.Model.World, &ui)

	view := res.NewView(48, 24)
	ecs.AddResource(&g.Model.World, &view)

	production := res.Production{}
	stock := res.Stock{}
	ecs.AddResource(&g.Model.World, &production)
	ecs.AddResource(&g.Model.World, &stock)

	// =========== Systems ===========

	g.Model.AddSystem(&sys.InitTerrain{})

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

	// =========== Run ===========

	g.Initialize()
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
