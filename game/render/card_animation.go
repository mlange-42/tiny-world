package render

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// CardAnimation is a system to render card animations.
type CardAnimation struct {
	MaxOffset int
	Duration  int

	time    generic.Resource[res.GameTick]
	screen  generic.Resource[res.Screen]
	sprites generic.Resource[res.Sprites]
	view    generic.Resource[res.View]

	filter generic.Filter1[comp.CardAnimation]

	toRemove []ecs.Entity
}

// InitializeUI the system
func (s *CardAnimation) InitializeUI(world *ecs.World) {
	s.time = generic.NewResource[res.GameTick](world)
	s.screen = generic.NewResource[res.Screen](world)
	s.sprites = generic.NewResource[res.Sprites](world)
	s.view = generic.NewResource[res.View](world)

	s.filter = *generic.NewFilter1[comp.CardAnimation]()
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

		if below != terr.Air {
			bIdx := sprites.GetTerrainIndex(below)
			sp := sprites.Get(bIdx)
			img.DrawImage(sp, &op)
		}

		idx := sprites.GetTerrainIndex(t)
		sp := sprites.GetRand(idx, 0, rnd)
		img.DrawImage(sp, &op)
	}

	query := s.filter.Query(world)
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
