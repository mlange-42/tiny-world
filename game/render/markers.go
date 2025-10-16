package render

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
)

// Markers is a system to render production markers.
type Markers struct {
	MinOffset int
	MaxOffset int
	Duration  int

	time    ecs.Resource[res.GameTick]
	screen  ecs.Resource[res.Screen]
	sprites ecs.Resource[res.Sprites]
	view    ecs.Resource[res.View]

	filter *ecs.Filter2[comp.Tile, comp.ProductionMarker]

	resources []int
}

// InitializeUI the system
func (s *Markers) InitializeUI(world *ecs.World) {
	s.time = ecs.NewResource[res.GameTick](world)
	s.screen = s.screen.New(world)
	s.sprites = ecs.NewResource[res.Sprites](world)
	s.view = ecs.NewResource[res.View](world)

	s.filter = s.filter.New(world)

	sprites := s.sprites.Get()
	s.resources = make([]int, len(resource.Properties))
	for i := range resource.Properties {
		s.resources[i] = sprites.GetIndex(resource.Properties[i].Name)
	}
}

// UpdateUI the system
func (s *Markers) UpdateUI(world *ecs.World) {
	tick := s.time.Get().Tick
	sprites := s.sprites.Get()
	view := s.view.Get()
	canvas := s.screen.Get()
	img := canvas.Image

	off := view.Offset()
	bounds := view.Bounds(canvas.Width, canvas.Height)

	op := ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver
	if view.Zoom < 1 {
		op.Filter = ebiten.FilterLinear
	}

	halfWidth := view.TileWidth / 2

	drawSprite := func(point *image.Point, cursor int) {
		info := sprites.GetInfo(cursor)
		sp := sprites.Get(cursor)
		h := sp.Bounds().Dy() - view.TileHeight

		op.GeoM.Reset()
		op.GeoM.Scale(view.Zoom, view.Zoom)
		op.GeoM.Translate(
			float64(point.X-halfWidth)*view.Zoom-float64(off.X),
			float64(point.Y-h-info.YOffset)*view.Zoom-float64(off.Y),
		)
		img.DrawImage(sp, &op)
	}

	query := s.filter.Query()
	for query.Next() {
		tile, mark := query.Get()
		point := view.TileToGlobal(tile.X, tile.Y)
		if !point.In(bounds) {
			continue
		}
		passed := tick - mark.StartTick
		off := s.MinOffset + (s.MaxOffset-s.MinOffset)*int(passed)/s.Duration
		point.Y -= off
		drawSprite(&point, s.resources[mark.Resource])
	}
}

// PostUpdateUI the system
func (s *Markers) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *Markers) FinalizeUI(world *ecs.World) {}
