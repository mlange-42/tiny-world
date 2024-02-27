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
	update   *res.UpdateInterval

	prodMapper   generic.Map1[comp.Production]
	pathMapper   generic.Map1[comp.Path]
	haulerMapper generic.Map2[comp.Hauler, comp.HaulerSprite]

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
	updateRes := generic.NewResource[res.UpdateInterval](world)

	s.rules = rulesRes.Get()
	s.view = viewRes.Get()
	s.sprites = spritesRes.Get()
	s.terrain = terrainRes.Get()
	s.landUse = landUseRes.Get()
	s.landUseE = landUseERes.Get()
	s.update = updateRes.Get()

	s.prodMapper = generic.NewMap1[comp.Production](world)
	s.pathMapper = generic.NewMap1[comp.Path](world)
	s.haulerMapper = generic.NewMap2[comp.Hauler, comp.HaulerSprite](world)

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
				height = s.drawSprite(img, &s.terrain.TerrainGrid, i, j, t, &point, height, &off, false)
			}

			lu := s.landUse.Get(i, j)
			if lu != terr.Air {
				if terr.Buildings.Contains(lu) {
					_ = s.drawSprite(img, &s.landUse.TerrainGrid, i, j, terr.Path, &point, height, &off, true)
				}
				_ = s.drawSprite(img, &s.landUse.TerrainGrid, i, j, lu, &point, height, &off, false)
			}

			if lu == terr.Path {
				path := s.pathMapper.Get(s.landUseE.Get(i, j))
				offset := 0.1
				for _, h := range path.Haulers {
					haul, sp := s.haulerMapper.Get(h.Entity)
					s.drawHauler(img, sp.SpriteIndex, haul, height, offset, &off)
				}
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

func (s *Terrain) drawHauler(img *ebiten.Image, sprite int, haul *comp.Hauler, height int, offset float64, camOffset *image.Point) {
	p1 := haul.Path[haul.Index]
	p2 := haul.Path[haul.Index-1]
	midX, midY := float64(p1.X+p2.X)/2, float64(p1.Y+p2.Y)/2

	dx, dy := float64(p2.X-p1.X), float64(p2.Y-p1.Y)
	dxr, dyr := -dy, dx

	var dxStart, dyStart float64
	var dxEnd, dyEnd float64
	var xx, yy float64

	dt := float64(haul.PathFraction) / float64(s.update.Interval)
	if dt <= 0.5 {
		if haul.Index < len(haul.Path)-1 {
			p3 := haul.Path[haul.Index+1]
			dx2, dy2 := float64(p1.X-p3.X), float64(p1.Y-p3.Y)
			if !(dx2 == dx && dy2 == dy) {
				if dx2 == dxr && dy2 == dyr {
					dxStart, dyStart = -dx, -dy
				} else {
					dxStart, dyStart = dx, dy
				}
			}
		}
		frac := dt * 2
		xx = (float64(p1.X)+dxStart*offset)*(1-frac) + midX*frac
		yy = (float64(p1.Y)+dyStart*offset)*(1-frac) + midY*frac
	} else {
		if haul.Index > 1 {
			p3 := haul.Path[haul.Index-2]
			dx2, dy2 := float64(p3.X-p2.X), float64(p3.Y-p2.Y)
			if !(dx2 == dx && dy2 == dy) {
				if dx2 == dxr && dy2 == dyr {
					dxEnd, dyEnd = dx, dy
				} else {
					dxEnd, dyEnd = -dx, -dy
				}
			}
		}
		frac := (dt - 0.5) * 2
		xx = midX*(1-frac) + (float64(p2.X)-dxEnd*offset)*frac
		yy = midY*(1-frac) + (float64(p2.Y)-dyEnd*offset)*frac
	}

	pt := s.view.SubtileToGlobal(xx+offset*dxr, yy+offset*dyr)

	s.drawSimpleSprite(img, sprite, &pt, height, camOffset)
}

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
			s.drawSprite(img, &s.terrain.TerrainGrid, x, y, toBuild, point, height, camOffset, true)
		} else {
			canBuildHere = canBuildHere && luEntity.IsZero()
			s.drawSprite(img, &s.landUse.TerrainGrid, x, y, toBuild, point, height, camOffset, true)
		}
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
	prod := s.prodMapper.Get(luEntity)
	text.Draw(img, fmt.Sprintf("%d (%d)", prod.Amount, prod.Stock), s.font,
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
		float64(point.X-sp.Bounds().Dx()/2)*z-float64(camOffset.X),
		float64(point.Y-h-info.YOffset)*z-float64(camOffset.Y),
	)
	img.DrawImage(sp, &op)
}

func (s *Terrain) drawSprite(img *ebiten.Image, grid *res.TerrainGrid,
	x, y int, t terr.Terrain, point *image.Point, height int,
	camOffset *image.Point,
	selfConnect bool) int {

	idx := s.sprites.GetTerrainIndex(t)
	sp, info := s.sprites.Get(idx)
	h := sp.Bounds().Dy() - s.view.TileHeight

	if info.IsMultitile() {
		var neigh terr.Directions
		if selfConnect {
			neigh = grid.NeighborsMask(x, y, t)
		} else {
			neigh = grid.NeighborsMaskMulti(x, y, terr.Properties[t].ConnectsTo)
		}
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
		float64(point.X-sp.Bounds().Dx()/2)*z-float64(camOffset.X),
		float64(point.Y-h-height-info.YOffset)*z-float64(camOffset.Y),
	)
	img.DrawImage(sp, &op)

	return height + info.Height
}

func (s *Terrain) drawSimpleSprite(img *ebiten.Image,
	idx int, point *image.Point, height int,
	camOffset *image.Point) int {

	sp, info := s.sprites.Get(idx)
	h := sp.Bounds().Dy() - s.view.TileHeight

	op := ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver
	if s.view.Zoom < 1 {
		op.Filter = ebiten.FilterLinear
	}

	z := s.view.Zoom
	op.GeoM.Scale(z, z)
	op.GeoM.Translate(
		float64(point.X-sp.Bounds().Dx()/2)*z-float64(camOffset.X),
		float64(point.Y-h-height-info.YOffset)*z-float64(camOffset.Y),
	)
	img.DrawImage(sp, &op)

	return height + info.Height
}
