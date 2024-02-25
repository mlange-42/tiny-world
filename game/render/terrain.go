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
	"golang.org/x/image/font"
)

// Terrain is a system to render the terrain.
type Terrain struct {
	cursorGreen  int
	cursorRed    int
	cursorBlue   int
	cursorYellow int

	screen    generic.Resource[res.EbitenImage]
	selection generic.Resource[res.Selection]

	rules    *res.Rules
	view     *res.View
	sprites  *res.Sprites
	terrain  *res.Terrain
	landUse  *res.LandUse
	landUseE *res.LandUseEntities

	prodMapper generic.Map1[comp.Production]

	font font.Face
}

// InitializeUI the system
func (s *Terrain) InitializeUI(world *ecs.World) {
	s.screen = generic.NewResource[res.EbitenImage](world)
	s.selection = generic.NewResource[res.Selection](world)

	rulesRes := generic.NewResource[res.Rules](world)
	viewRes := generic.NewResource[res.View](world)
	spritesRes := generic.NewResource[res.Sprites](world)
	terrainRes := generic.NewResource[res.Terrain](world)
	landUseRes := generic.NewResource[res.LandUse](world)
	landUseERes := generic.NewResource[res.LandUseEntities](world)

	s.rules = rulesRes.Get()
	s.view = viewRes.Get()
	s.sprites = spritesRes.Get()
	s.terrain = terrainRes.Get()
	s.landUse = landUseRes.Get()
	s.landUseE = landUseERes.Get()

	s.prodMapper = generic.NewMap1[comp.Production](world)

	s.cursorRed = s.sprites.GetIndex("cursor_red")
	s.cursorGreen = s.sprites.GetIndex("cursor_green")
	s.cursorBlue = s.sprites.GetIndex("cursor_blue")
	s.cursorYellow = s.sprites.GetIndex("cursor_yellow")

	fts := generic.NewResource[res.Fonts](world)
	fonts := fts.Get()
	s.font = fonts.Default
}

// UpdateUI the system
func (s *Terrain) UpdateUI(world *ecs.World) {
	sel := s.selection.Get()

	canvas := s.screen.Get()
	img := canvas.Image

	off := s.view.Offset()
	bounds := s.view.Bounds(canvas.Width, canvas.Height)

	img.Clear()

	mx, my := s.view.ScreenToGlobal(ebiten.CursorPosition())
	cursor := s.view.GlobalToTile(mx, my)

	for i := 0; i < s.terrain.Width(); i++ {
		for j := 0; j < s.terrain.Height(); j++ {
			point := s.view.TileToGlobal(i, j)
			if !point.In(bounds) {
				continue
			}

			height := 0
			t := s.terrain.Get(i, j)
			if t != terr.Air && t != terr.Buildable {
				height = s.drawSprite(img, &s.terrain.Grid, i, j, t, &point, height, &off)
			}

			lu := s.landUse.Get(i, j)
			if lu != terr.Air {
				_ = s.drawSprite(img, &s.landUse.Grid, i, j, lu, &point, height, &off)
			}

			if cursor.X == i && cursor.Y == j {
				s.drawCursor(img,
					i, j, height, &point, &off, sel.BuildType)
			}
		}
	}
}

// PostUpdateUI the system
func (s *Terrain) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *Terrain) FinalizeUI(world *ecs.World) {}

func (s *Terrain) drawCursor(img *ebiten.Image,
	x, y, height int, point *image.Point, camOffset *image.Point,
	toBuild terr.Terrain) {

	ter := s.terrain.Get(x, y)
	lu := s.landUse.Get(x, y)
	luEntity := s.landUseE.Get(x, y)
	prop := terr.Properties[toBuild]
	if prop.CanBuild {
		canDestroy := lu == toBuild &&
			(s.rules.AllowRemoveBuilt && prop.CanBuy) || (s.rules.AllowRemoveNatural && !prop.CanBuy)

		canBuildHere := prop.BuildOn.Contains(ter)
		if prop.IsTerrain {
			height = 0
		} else {
			canBuildHere = canBuildHere && luEntity.IsZero()
		}
		s.drawSprite(img, &s.landUse.Grid, x, y, toBuild, point, height, camOffset)
		if canBuildHere {
			s.drawCursorSprite(img, point, camOffset, s.cursorGreen)
		} else {
			if canDestroy {
				s.drawCursorSprite(img, point, camOffset, s.cursorYellow)
			} else {
				s.drawCursorSprite(img, point, camOffset, s.cursorRed)
			}
		}
	} else {
		s.drawCursorSprite(img, point, camOffset, s.cursorBlue)
	}

	if terr.Properties[lu].Production.Produces == resource.EndResources {
		return
	}

	if luEntity.IsZero() {
		return
	}
	prod := s.prodMapper.Get(luEntity).Amount
	text.Draw(img, fmt.Sprint(prod), s.font,
		int(float64(point.X-s.view.TileWidth/2)*s.view.Zoom-float64(camOffset.X)),
		int(float64(point.Y-2*s.view.TileHeight)*s.view.Zoom-float64(camOffset.Y)),
		color.RGBA{255, 255, 255, 255},
	)
}
func (s *Terrain) drawCursorSprite(img *ebiten.Image,
	point *image.Point, camOffset *image.Point, cursor int) {

	op := ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver
	if s.view.Zoom < 1 {
		op.Filter = ebiten.FilterLinear
	}

	sp, info := s.sprites.Get(cursor)
	h := sp.Bounds().Dy() - s.view.TileHeight

	z := s.view.Zoom
	op.GeoM.Scale(z, z)
	op.GeoM.Translate(
		float64(point.X-s.view.TileWidth/2)*z-float64(camOffset.X),
		float64(point.Y-h-info.YOffset)*z-float64(camOffset.Y),
	)
	img.DrawImage(sp, &op)
}

func (s *Terrain) drawSprite(img *ebiten.Image, grid *res.Grid[terr.Terrain],
	x, y int, t terr.Terrain, point *image.Point, height int,
	camOffset *image.Point) int {

	idx := s.sprites.GetTerrainIndex(t)
	sp, info := s.sprites.Get(idx)
	h := sp.Bounds().Dy() - s.view.TileHeight

	if info.MultiTile {
		neigh := grid.NeighborsMask(x, y, t)
		idx = s.sprites.GetMultiTileIndex(t, neigh)
		sp, _ = s.sprites.Get(idx)
	}

	op := ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver
	if s.view.Zoom < 1 {
		op.Filter = ebiten.FilterLinear
	}

	z := s.view.Zoom
	op.GeoM.Scale(z, z)
	op.GeoM.Translate(
		float64(point.X-s.view.TileWidth/2)*z-float64(camOffset.X),
		float64(point.Y-h-height-info.YOffset)*z-float64(camOffset.Y),
	)
	img.DrawImage(sp, &op)

	return height + info.Height
}
