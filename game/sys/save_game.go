package sys

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
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

	ui        ecs.Resource[res.UI]
	saveEvent ecs.Resource[res.SaveEvent]
	saveTime  ecs.Resource[res.SaveTime]
	skip      []ecs.Comp
}

// Initialize the system
func (s *SaveGame) Initialize(world *ecs.World) {
	s.ui = ecs.NewResource[res.UI](world)
	s.saveEvent = ecs.NewResource[res.SaveEvent](world)
	s.saveTime = ecs.NewResource[res.SaveTime](world)

	s.skip = []ecs.Comp{ecs.C[res.Fonts](),
		ecs.C[res.Screen](),
		ecs.C[res.EntityFactory](),
		ecs.C[res.Sprites](),
		ecs.C[res.Terrain](),
		ecs.C[res.TerrainEntities](),
		ecs.C[res.LandUse](),
		ecs.C[res.LandUseEntities](),
		ecs.C[res.Buildable](),
		ecs.C[res.SaveEvent](),
		ecs.C[res.UpdateInterval](),
		ecs.C[res.GameSpeed](),
		ecs.C[res.Selection](),
		ecs.C[achievements.Achievements](),
		ecs.C[res.Mouse](),
		ecs.C[res.View](),
		ecs.C[res.UI](),
		ecs.C[resource.Termination](),
		ecs.C[resource.Rand](),
		ecs.C[app.Systems](),
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

	if evt.ShouldSave || (keysPressed && !shift) {
		evt.ShouldSave = false
		print("Saving game... ")

		s.saveTime.Get().Time = time.Now()
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
