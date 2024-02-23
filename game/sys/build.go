package sys

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Build system.
type Build struct {
	view      generic.Resource[res.View]
	terrain   generic.Resource[res.Terrain]
	landUse   generic.Resource[res.LandUse]
	selection generic.Resource[res.Selection]
}

// Initialize the system
func (s *Build) Initialize(world *ecs.World) {
	s.view = generic.NewResource[res.View](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.selection = generic.NewResource[res.Selection](world)
}

// Update the system
func (s *Build) Update(world *ecs.World) {
	sel := s.selection.Get()

	for i := range terr.Properties {
		p := &terr.Properties[i]
		if p.CanBuild && inpututil.IsKeyJustPressed(p.ShortKey) {
			sel.Build = terr.Terrain(i)
			fmt.Println("paint", p.Name)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sel.Build = terr.Air
		fmt.Println("paint nothing")
	}

	p := &terr.Properties[sel.Build]
	if !p.CanBuild ||
		!(inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2)) {
		return
	}

	view := s.view.Get()
	landUse := s.landUse.Get()
	mx, my := view.ScreenToGlobal(ebiten.CursorPosition())
	cursor := view.GlobalToTile(mx, my)

	remove := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2)
	if remove {
		if p.IsTerrain {
			return
		}
		luHere := landUse.Get(cursor.X, cursor.Y)
		if luHere == sel.Build {
			landUse.Set(cursor.X, cursor.Y, terr.Air)
		}
		return
	}

	terrain := s.terrain.Get()

	terrHere := terrain.Get(cursor.X, cursor.Y)
	if !p.BuildOn.Contains(terrHere) {
		return
	}
	if p.IsTerrain {
		terrain.Set(cursor.X, cursor.Y, sel.Build)
	} else {
		landUse.Set(cursor.X, cursor.Y, sel.Build)
	}
}

// Finalize the system
func (s *Build) Finalize(world *ecs.World) {}
