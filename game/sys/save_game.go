package sys

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/resource"
	as "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// SaveGame system.
type SaveGame struct {
	Path string
}

// Initialize the system
func (s *SaveGame) Initialize(world *ecs.World) {}

// Update the system
func (s *SaveGame) Update(world *ecs.World) {
	if !(ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyS)) {
		return
	}
	print("Saving game... ")

	js, err := as.Serialize(world,
		as.Opts.SkipResources(
			generic.T[res.Fonts](),
			generic.T[res.EbitenImage](),
			generic.T[res.EntityFactory](),
			generic.T[res.Sprites](),
			generic.T[res.Terrain](),
			generic.T[res.TerrainEntities](),
			generic.T[res.LandUse](),
			generic.T[res.LandUseEntities](),
			generic.T[resource.Termination](),
			generic.T[model.Systems](),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(s.Path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Printf("Error saving game: %s", err.Error())
		return
	}

	f, err := os.Create(s.Path)
	if err != nil {
		log.Printf("Error saving game: %s", err.Error())
		return
	}
	defer f.Close()

	f.Write(js)
	println("done.")
}

// Finalize the system
func (s *SaveGame) Finalize(world *ecs.World) {}
