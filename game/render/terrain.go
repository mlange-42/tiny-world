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

	time     *res.GameTick
	rules    *res.Rules
	view     *res.View
	sprites  *res.Sprites
	terrain  *res.Terrain
	terrainE *res.TerrainEntities
	landUse  *res.LandUse
	landUseE *res.LandUseEntities
	update   *res.UpdateInterval

	prodMapper   generic.Map2[comp.Terrain, comp.Production]
	pathMapper   generic.Map1[comp.Path]
	haulerMapper generic.Map2[comp.Hauler, comp.HaulerSprite]
	spriteMapper generic.Map1[comp.RandomSprite]

	font font.Face
}

// InitializeUI the system
func (s *Terrain) InitializeUI(world *ecs.World) {
	s.screen = generic.NewResource[res.EbitenImage](world)
	s.selection = generic.NewResource[res.Selection](world)

	s.time = ecs.GetResource[res.GameTick](world)
	s.rules = ecs.GetResource[res.Rules](world)
	s.view = ecs.GetResource[res.View](world)
	s.sprites = ecs.GetResource[res.Sprites](world)
	s.terrain = ecs.GetResource[res.Terrain](world)
	s.terrainE = ecs.GetResource[res.TerrainEntities](world)
	s.landUse = ecs.GetResource[res.LandUse](world)
	s.landUseE = ecs.GetResource[res.LandUseEntities](world)
	s.update = ecs.GetResource[res.UpdateInterval](world)

	s.prodMapper = generic.NewMap2[comp.Terrain, comp.Production](world)
	s.pathMapper = generic.NewMap1[comp.Path](world)
	s.haulerMapper = generic.NewMap2[comp.Hauler, comp.HaulerSprite](world)
	s.spriteMapper = generic.NewMap1[comp.RandomSprite](world)

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
	img.Fill(s.sprites.Background)

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
				tE := s.terrainE.Get(i, j)
				randTile := s.spriteMapper.Get(tE)
				height = s.drawSprite(img, s.terrain, s.landUse, i, j, t, &point, height, &off, randTile, false, terr.Properties[t].TerrainBelow)
			}

			lu := s.landUse.Get(i, j)
			if lu != terr.Air {
				luE := s.landUseE.Get(i, j)
				randTile := s.spriteMapper.Get(luE)
				_ = s.drawSprite(img, s.terrain, s.landUse, i, j, lu, &point, height, &off, randTile, false, terr.Properties[lu].TerrainBelow)
			}

			if terr.Properties[lu].IsPath {
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
		canDestroy := lu == toBuild && prop.CanBuy

		canBuildHere := prop.BuildOn.Contains(ter)
		if prop.IsTerrain {
			height = 0
		} else {
			luNatural := !terr.Properties[lu].CanBuy
			canBuildHere = canBuildHere && (lu == terr.Air || (luNatural && prop.CanBuy))
		}
		s.drawSprite(img, s.terrain, s.landUse, x, y, toBuild, point, height, camOffset, nil, true, prop.TerrainBelow)

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

	propHere := &terr.Properties[lu]
	if propHere.Production.MaxProduction == 0 {
		return
	}

	if luEntity.IsZero() {
		return
	}
	tp, prod := s.prodMapper.Get(luEntity)
	bStock := terr.Properties[tp.Terrain].Storage[prod.Resource]
	text.Draw(img, fmt.Sprintf("%d/%d (%d/%d)", prod.Amount, propHere.Production.MaxProduction, prod.Stock, bStock), s.font,
		int(float64(point.X)*s.view.Zoom-32-float64(camOffset.X)),
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

	info := s.sprites.GetInfo(cursor)
	sp := s.sprites.Get(cursor)

	h := sp.Bounds().Dy() - s.view.TileHeight

	z := s.view.Zoom
	op.GeoM.Scale(z, z)
	op.GeoM.Translate(
		float64(point.X-sp.Bounds().Dx()/2)*z-float64(camOffset.X),
		float64(point.Y-h-info.YOffset)*z-float64(camOffset.Y),
	)
	img.DrawImage(sp, &op)
}

func (s *Terrain) drawSprite(img *ebiten.Image, terrain *res.Terrain, landUse *res.LandUse,
	x, y int, t terr.Terrain, point *image.Point, height int,
	camOffset *image.Point, randSprite *comp.RandomSprite,
	selfConnect bool, below terr.Terrain) int {

	idx := s.sprites.GetTerrainIndex(t)
	info := s.sprites.GetInfo(idx)

	if below != terr.Air {
		height = s.drawSprite(img, terrain, landUse,
			x, y, below, point, height,
			camOffset, randSprite,
			terr.Properties[t].SelfConnectBelow, terr.Air)
	}

	var sp *ebiten.Image
	if info.IsMultitile() {
		var neigh terr.Directions
		if selfConnect {
			neigh = terrain.NeighborsMask(x, y, t) | landUse.NeighborsMask(x, y, t)
		} else {
			conn := terr.Properties[t].ConnectsTo
			neigh = terrain.NeighborsMaskMulti(x, y, conn) | landUse.NeighborsMaskMulti(x, y, conn)
		}
		mIdx := s.sprites.GetMultiTileIndex(t, neigh, int(s.time.Tick), int(randSprite.GetRand()))

		sp = s.sprites.GetSprite(mIdx)
	} else {
		sp = s.sprites.GetRand(idx, int(s.time.Tick), int(randSprite.GetRand()))
	}
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

func (s *Terrain) drawSimpleSprite(img *ebiten.Image,
	idx int, point *image.Point, height int,
	camOffset *image.Point) int {

	info := s.sprites.GetInfo(idx)
	sp := s.sprites.Get(idx)
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
