package sys

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Build system.
type Build struct {
	paint terr.Terrain

	view    generic.Resource[game.View]
	terrain generic.Resource[game.Terrain]
	landUse generic.Resource[game.LandUse]
}

// Initialize the system
func (s *Build) Initialize(world *ecs.World) {
	s.view = generic.NewResource[game.View](world)
	s.terrain = generic.NewResource[game.Terrain](world)
	s.landUse = generic.NewResource[game.LandUse](world)
}

// Update the system
func (s *Build) Update(world *ecs.World) {
	for i := range terr.Properties {
		p := &terr.Properties[i]
		if p.CanBuild && inpututil.IsKeyJustPressed(p.ShortKey) {
			s.paint = terr.Terrain(i)
			fmt.Println("paint", p.Name)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.paint = terr.EndTerrain
		fmt.Println("paint nothing")
	}

	p := &terr.Properties[s.paint]
	if s.paint == terr.EndTerrain || !p.CanBuild || !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		return
	}

	view := s.view.Get()
	terrain := s.terrain.Get()
	landUse := s.landUse.Get()

	mx, my := view.ScreenToGlobal(ebiten.CursorPosition())
	cursor := view.GlobalToTile(mx, my)

	terrHere := terrain.Get(cursor.X, cursor.Y)
	if !p.BuildOn.Contains(terrHere) {
		return
	}

	if p.IsTerrain {
		terrain.Set(cursor.X, cursor.Y, s.paint)
	} else {
		landUse.Set(cursor.X, cursor.Y, s.paint)
	}
}

// Finalize the system
func (s *Build) Finalize(world *ecs.World) {}
