package sys

import (
	"fmt"
	"log"

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
type SaveGame struct{}

// Initialize the system
func (s *SaveGame) Initialize(world *ecs.World) {}

// Update the system
func (s *SaveGame) Update(world *ecs.World) {
	if !(ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyS)) {
		return
	}
	println("Saving game")

	js, err := as.Serialize(world,
		as.Opts.SkipResources(
			generic.T[res.Fonts](),
			generic.T[res.HUD](),
			generic.T[res.EbitenImage](),
			generic.T[res.Sprites](),
			generic.T[res.UI](),
			generic.T[resource.Termination](),
			generic.T[model.Systems](),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(js))
}

// Finalize the system
func (s *SaveGame) Finalize(world *ecs.World) {}
