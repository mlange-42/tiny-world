package sys

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/save"
)

// SaveGame system.
type SaveGame struct {
	Folder string
	Name   string

	saveEvent generic.Resource[res.SaveEvent]
}

// Initialize the system
func (s *SaveGame) Initialize(world *ecs.World) {
	s.saveEvent = generic.NewResource[res.SaveEvent](world)
}

// Update the system
func (s *SaveGame) Update(world *ecs.World) {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		err := save.SaveMap(world)
		if err != nil {
			log.Printf("Error saving map: %s", err.Error())
			return
		}
		return
	}

	if !((ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyS)) ||
		s.saveEvent.Get().ShouldSave) {
		return
	}
	s.saveEvent.Get().ShouldSave = false
	print("Saving game... ")

	err := save.SaveWorld(s.Folder, s.Name, world)
	if err != nil {
		log.Printf("Error saving game: %s", err.Error())
		return
	}
	println("done.")
}

// Finalize the system
func (s *SaveGame) Finalize(world *ecs.World) {}
