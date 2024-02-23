package render

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Terrain is a system to render the terrain.
type Terrain struct {
	cursorGreen int
	cursorRed   int
	cursorBlue  int

	screen    generic.Resource[res.EbitenImage]
	sprites   generic.Resource[res.Sprites]
	terrain   generic.Resource[res.Terrain]
	landUse   generic.Resource[res.LandUse]
	landUseE  generic.Resource[res.LandUseEntities]
	view      generic.Resource[res.View]
	selection generic.Resource[res.Selection]
	fonts     generic.Resource[res.Fonts]

	prodMapper generic.Map1[comp.Production]
}

// InitializeUI the system
func (s *Terrain) InitializeUI(world *ecs.World) {
	s.screen = generic.NewResource[res.EbitenImage](world)
	s.sprites = generic.NewResource[res.Sprites](world)
	s.terrain = generic.NewResource[res.Terrain](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.landUseE = generic.NewResource[res.LandUseEntities](world)
	s.view = generic.NewResource[res.View](world)
	s.selection = generic.NewResource[res.Selection](world)
	s.fonts = generic.NewResource[res.Fonts](world)

	s.prodMapper = generic.NewMap1[comp.Production](world)

	sprites := s.sprites.Get()
	s.cursorRed = sprites.GetIndex("cursor_red")
	s.cursorGreen = sprites.GetIndex("cursor_green")
	s.cursorBlue = sprites.GetIndex("cursor_blue")
}

// UpdateUI the system
func (s *Terrain) UpdateUI(world *ecs.World) {
	terrain := s.terrain.Get()
	landUse := s.landUse.Get()
	landUseE := s.landUseE.Get()
	sprites := s.sprites.Get()
	view := s.view.Get()
	sel := s.selection.Get()
	fonts := s.fonts.Get()

	canvas := s.screen.Get()
	img := canvas.Image

	off := view.Offset()
	bounds := view.Bounds(canvas.Width, canvas.Height)

	img.Clear()

	op := ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver
	if view.Zoom < 1 {
		op.Filter = ebiten.FilterLinear
	}

	halfWidth := view.TileWidth / 2

	drawSprite := func(grid *res.Grid[terr.Terrain], x, y int, t terr.Terrain, point *image.Point, height int) int {
		idx := sprites.GetTerrainIndex(t)
		sp, info := sprites.Get(idx)
		h := sp.Bounds().Dy() - view.TileHeight

		if info.MultiTile {
			neigh := grid.NeighborsMask(x, y, t)
			idx = sprites.GetMultiTileIndex(t, neigh)
			sp, _ = sprites.Get(idx)
		}

		op.GeoM.Reset()
		op.GeoM.Scale(view.Zoom, view.Zoom)
		op.GeoM.Translate(
			float64(point.X-halfWidth)*view.Zoom-float64(off.X),
			float64(point.Y-h-height-info.YOffset)*view.Zoom-float64(off.Y),
		)
		img.DrawImage(sp, &op)

		return height + info.Height
	}

	drawCursor := func(point *image.Point, cursor int) {
		sp, info := sprites.Get(cursor)
		h := sp.Bounds().Dy() - view.TileHeight

		op.GeoM.Reset()
		op.GeoM.Scale(view.Zoom, view.Zoom)
		op.GeoM.Translate(
			float64(point.X-halfWidth)*view.Zoom-float64(off.X),
			float64(point.Y-h-info.YOffset)*view.Zoom-float64(off.Y),
		)
		img.DrawImage(sp, &op)
	}

	mx, my := view.ScreenToGlobal(ebiten.CursorPosition())
	cursor := view.GlobalToTile(mx, my)

	for i := 0; i < terrain.Width(); i++ {
		for j := 0; j < terrain.Height(); j++ {
			point := view.TileToGlobal(i, j)
			if !point.In(bounds) {
				continue
			}

			height := 0
			t := terrain.Get(i, j)
			if t != terr.Air && t != terr.Buildable {
				height = drawSprite(&terrain.Grid, i, j, t, &point, height)
			}

			lu := landUse.Get(i, j)
			if lu != terr.Air {
				_ = drawSprite(&landUse.Grid, i, j, lu, &point, height)
			}

			if cursor.X == i && cursor.Y == j {
				prop := terr.Properties[sel.Build]
				if prop.CanBuild {
					terrHere := terrain.Get(cursor.X, cursor.Y)
					if prop.BuildOn.Contains(terrHere) {
						drawCursor(&point, s.cursorGreen)
					} else {
						drawCursor(&point, s.cursorRed)
					}
				} else {
					drawCursor(&point, s.cursorBlue)
				}

				lu := landUse.Get(i, j)
				if terr.Properties[lu].Production.Produces == resource.EndResources {
					continue
				}

				luEntity := landUseE.Get(i, j)
				if luEntity.IsZero() {
					continue
				}
				prod := s.prodMapper.Get(luEntity).Amount
				text.Draw(img, fmt.Sprint(prod), fonts.Default,
					int(float64(point.X-halfWidth/2)*view.Zoom-float64(off.X)),
					int(float64(point.Y-2*view.TileHeight)*view.Zoom-float64(off.Y)),
					color.RGBA{255, 255, 255, 255},
				)
			}
		}
	}
}

// PostUpdateUI the system
func (s *Terrain) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *Terrain) FinalizeUI(world *ecs.World) {}
