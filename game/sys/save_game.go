package sys

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/save"
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

	err := save.SaveWorld(s.Path, world)
	if err != nil {
		log.Printf("Error saving game: %s", err.Error())
		return
	}
	println("done.")
}

// Finalize the system
func (s *SaveGame) Finalize(world *ecs.World) {}
