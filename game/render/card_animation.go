package render

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// CardAnimation is a system to render card animations.
type CardAnimation struct {
	MaxOffset int
	Duration  int

	time    ecs.Resource[res.GameTick]
	screen  ecs.Resource[res.Screen]
	sprites ecs.Resource[res.Sprites]
	view    ecs.Resource[res.View]

	filter *ecs.Filter1[comp.CardAnimation]

	toRemove []ecs.Entity
}

// InitializeUI the system
func (s *CardAnimation) InitializeUI(world *ecs.World) {
	s.time = s.time.New(world)
	s.screen = s.screen.New(world)
	s.sprites = s.sprites.New(world)
	s.view = s.view.New(world)

	s.filter = s.filter.New(world)
}

// UpdateUI the system
func (s *CardAnimation) UpdateUI(world *ecs.World) {
	tick := s.time.Get().RenderTick
	sprites := s.sprites.Get()
	canvas := s.screen.Get()
	img := canvas.Image

	op := ebiten.DrawImageOptions{}

	drawSprite := func(point *image.Point, t terr.Terrain, rnd int) {
		below := terr.Properties[t].TerrainBelow

		op.GeoM.Reset()
		op.GeoM.Translate(float64(point.X), float64(point.Y))

		for _, tr := range below {
			bIdx := sprites.GetTerrainIndex(tr)
			sp := sprites.Get(bIdx)
			img.DrawImage(sp, &op)
		}

		idx := sprites.GetTerrainIndex(t)
		sp := sprites.GetRand(idx, 0, rnd)
		img.DrawImage(sp, &op)
	}

	query := s.filter.Query()
	for query.Next() {
		card := query.Get()

		passed := tick - card.StartTick
		if passed > int64(s.Duration) {
			s.toRemove = append(s.toRemove, query.Entity())
			continue
		}
		off := s.MaxOffset * int(passed) / s.Duration
		diff := card.Target.Sub(image.Pt(sprites.TileWidth/2, sprites.TileHeight/2)).Sub(card.Point)
		ln := 1 / math.Sqrt(float64(diff.X*diff.X+diff.Y*diff.Y))

		point := image.Pt(card.X+int(float64(diff.X*off)*ln), card.Y+int(float64(diff.Y*off)*ln))
		drawSprite(&point, card.Terrain, int(card.RandSprite))
	}

	for _, e := range s.toRemove {
		world.RemoveEntity(e)
	}
	s.toRemove = s.toRemove[:0]
}

// PostUpdateUI the system
func (s *CardAnimation) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *CardAnimation) FinalizeUI(world *ecs.World) {}
