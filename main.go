package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game"
	"github.com/mlange-42/tiny-world/game/render"
	"github.com/mlange-42/tiny-world/game/sys"
	"github.com/mlange-42/tiny-world/game/terr"
)

func main() {
	ebiten.SetVsyncEnabled(true)

	g := game.NewGame(model.New())

	ecs.AddResource(&g.Model.World, &g.Screen)

	terrain := game.Terrain{Grid: game.NewGrid[terr.Terrain](100, 100)}
	ecs.AddResource(&g.Model.World, &terrain)
	landUse := game.LandUse{Grid: game.NewGrid[terr.Terrain](100, 100)}
	ecs.AddResource(&g.Model.World, &landUse)

	view := game.View{
		TileWidth:   48,
		TileHeight:  24,
		Zoom:        1,
		MouseOffset: 24,
	}
	ecs.AddResource(&g.Model.World, &view)

	sprites := game.NewSprites("./assets/sprites")
	ecs.AddResource(&g.Model.World, &sprites)

	g.Model.AddSystem(&sys.InitTerrain{})
	g.Model.AddSystem(&sys.Build{})
	g.Model.AddSystem(&sys.PanAndZoom{
		PanButton: ebiten.MouseButton1,
	})

	g.Model.AddUISystem(&render.Terrain{})

	g.Initialize()
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
