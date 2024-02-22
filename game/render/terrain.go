package render

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Terrain is a system to render the terrain.
type Terrain struct {
	screen  generic.Resource[game.EbitenImage]
	sprites generic.Resource[game.Sprites]
	terrain generic.Resource[game.Terrain]
	view    generic.Resource[game.View]
}

// InitializeUI the system
func (s *Terrain) InitializeUI(world *ecs.World) {
	s.screen = generic.NewResource[game.EbitenImage](world)
	s.sprites = generic.NewResource[game.Sprites](world)
	s.terrain = generic.NewResource[game.Terrain](world)
	s.view = generic.NewResource[game.View](world)
}

// UpdateUI the system
func (s *Terrain) UpdateUI(world *ecs.World) {
	terrain := s.terrain.Get()
	sprites := s.sprites.Get()
	view := s.view.Get()

	canvas := s.screen.Get()
	img := canvas.Image

	off := view.Offset()
	bounds := view.Bounds(canvas.Width, canvas.Height)

	img.Clear()

	op := ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver

	drawSprite := func(idx int, point *image.Point, off *image.Point, height int) int {
		sp, info := sprites.Get(idx)
		h := sp.Bounds().Dy()

		op.GeoM.Reset()
		op.GeoM.Scale(view.Zoom, view.Zoom)
		op.GeoM.Translate(
			float64(point.X)*view.Zoom-float64(off.X),
			float64(point.Y-h-height-info.YOffset)*view.Zoom-float64(off.Y),
		)
		img.DrawImage(sp, &op)

		return height + info.Height
	}

	for i := 0; i < terrain.Width(); i++ {
		for j := 0; j < terrain.Height(); j++ {
			point := view.TileToScreen(i, j)
			if !point.In(bounds) {
				continue
			}

			t := terrain.Get(i, j)
			idx := sprites.GetTerrainIndex(t)

			height := 0
			if t != terr.Air {
				_ = drawSprite(idx, &point, &off, height)
			}
		}
	}
}

// PostUpdateUI the system
func (s *Terrain) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *Terrain) FinalizeUI(world *ecs.World) {}
