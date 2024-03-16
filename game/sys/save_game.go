package sys

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/resource"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/res/achievements"
	"github.com/mlange-42/tiny-world/game/save"
)

// SaveGame system.
type SaveGame struct {
	SaveFolder string
	MapFolder  string
	Name       string

	MainMenuFunc func()

	ui        generic.Resource[res.UI]
	saveEvent generic.Resource[res.SaveEvent]
	skip      []generic.Comp
}

// Initialize the system
func (s *SaveGame) Initialize(world *ecs.World) {
	s.ui = generic.NewResource[res.UI](world)
	s.saveEvent = generic.NewResource[res.SaveEvent](world)

	s.skip = []generic.Comp{generic.T[res.Fonts](),
		generic.T[res.Screen](),
		generic.T[res.EntityFactory](),
		generic.T[res.Sprites](),
		generic.T[res.Terrain](),
		generic.T[res.TerrainEntities](),
		generic.T[res.LandUse](),
		generic.T[res.LandUseEntities](),
		generic.T[res.Buildable](),
		generic.T[res.SaveEvent](),
		generic.T[res.UpdateInterval](),
		generic.T[res.GameSpeed](),
		generic.T[res.Selection](),
		generic.T[achievements.Achievements](),
		generic.T[res.Mouse](),
		generic.T[res.View](),
		generic.T[res.UI](),
		generic.T[resource.Termination](),
		generic.T[resource.Rand](),
		generic.T[model.Systems](),
	}
}

// Update the system
func (s *SaveGame) Update(world *ecs.World) {
	evt := s.saveEvent.Get()

	keysPressed := ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyS)
	shift := ebiten.IsKeyPressed(ebiten.KeyShift)

	if evt.ShouldSaveMap || (keysPressed && shift) {
		evt.ShouldSaveMap = false
		print("Saving map... ")
		err := save.SaveMap(s.MapFolder, s.Name, world)
		if err != nil {
			s.ui.Get().SetStatusLabel("Error saving map")
			log.Printf("Error saving map: %s", err.Error())
			return
		}
		s.ui.Get().SetStatusLabel("Map saved.")
		println("done.")
	}

	if evt.ShouldSaveMap || (keysPressed && !shift) {
		evt.ShouldSave = false
		print("Saving game... ")

		err := save.SaveWorld(s.SaveFolder, s.Name, world, s.skip)
		if err != nil {
			s.ui.Get().SetStatusLabel("Error saving game")
			log.Printf("Error saving game: %s", err.Error())
			return
		}
		s.ui.Get().SetStatusLabel("Game saved.")
		println("done.")
	}

	if evt.ShouldQuit {
		s.MainMenuFunc()
	}
}

// Finalize the system
func (s *SaveGame) Finalize(world *ecs.World) {}
