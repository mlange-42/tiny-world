package render

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/math"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/sprites"
	"github.com/mlange-42/tiny-world/game/terr"
	"golang.org/x/image/font"
)

// Terrain is a system to render the terrain.
type Terrain struct {
	cursorOk                    int
	cursorDenied                int
	cursorNeutral               int
	cursorDestroy               int
	warningMarker               int
	borderInner                 int
	borderOuter                 int
	indicatorPopulation         int
	indicatorPopulationInactive int
	indicatorProduction         int
	indicatorProductionInactive int
	indicatorStorage            int
	indicatorStorageInactive    int

	screen    generic.Resource[res.Screen]
	selection generic.Resource[res.Selection]
	mouse     generic.Resource[res.Mouse]
	ui        generic.Resource[res.UI]

	radiusFilter generic.Filter2[comp.Tile, comp.BuildRadius]

	time      *res.GameTick
	rules     *res.Rules
	view      *res.View
	sprites   *res.Sprites
	terrain   *res.Terrain
	terrainE  *res.TerrainEntities
	landUse   *res.LandUse
	landUseE  *res.LandUseEntities
	buildable *res.Buildable
	update    *res.UpdateInterval

	prodMapper    generic.Map2[comp.Terrain, comp.Production]
	popMapper     generic.Map1[comp.PopulationSupport]
	pathMapper    generic.Map1[comp.Path]
	haulerMapper  generic.Map2[comp.Hauler, comp.HaulerSprite]
	spriteMapper  generic.Map1[comp.RandomSprite]
	landUseMapper generic.Map2[comp.Production, comp.RandomSprite]

	font font.Face
}

// InitializeUI the system
func (s *Terrain) InitializeUI(world *ecs.World) {
	s.screen = generic.NewResource[res.Screen](world)
	s.selection = generic.NewResource[res.Selection](world)
	s.mouse = generic.NewResource[res.Mouse](world)
	s.ui = generic.NewResource[res.UI](world)

	s.time = ecs.GetResource[res.GameTick](world)
	s.rules = ecs.GetResource[res.Rules](world)
	s.view = ecs.GetResource[res.View](world)
	s.sprites = ecs.GetResource[res.Sprites](world)
	s.terrain = ecs.GetResource[res.Terrain](world)
	s.terrainE = ecs.GetResource[res.TerrainEntities](world)
	s.landUse = ecs.GetResource[res.LandUse](world)
	s.landUseE = ecs.GetResource[res.LandUseEntities](world)
	s.buildable = ecs.GetResource[res.Buildable](world)
	s.update = ecs.GetResource[res.UpdateInterval](world)

	s.prodMapper = generic.NewMap2[comp.Terrain, comp.Production](world)
	s.popMapper = generic.NewMap1[comp.PopulationSupport](world)
	s.pathMapper = generic.NewMap1[comp.Path](world)
	s.haulerMapper = generic.NewMap2[comp.Hauler, comp.HaulerSprite](world)
	s.spriteMapper = generic.NewMap1[comp.RandomSprite](world)
	s.landUseMapper = generic.NewMap2[comp.Production, comp.RandomSprite](world)

	s.radiusFilter = *generic.NewFilter2[comp.Tile, comp.BuildRadius]()

	s.cursorDenied = s.sprites.GetIndex(sprites.CursorDenied)
	s.cursorOk = s.sprites.GetIndex(sprites.CursorOk)
	s.cursorNeutral = s.sprites.GetIndex(sprites.CursorNeutral)
	s.cursorDestroy = s.sprites.GetIndex(sprites.CursorDestroy)
	s.warningMarker = s.sprites.GetIndex(sprites.WarningMarker)
	s.borderInner = s.sprites.GetIndex(sprites.BorderInner)
	s.borderOuter = s.sprites.GetIndex(sprites.BorderOuter)
	s.indicatorPopulation = s.sprites.GetIndex(sprites.IndicatorPopulation)
	s.indicatorPopulationInactive = s.sprites.GetIndex(sprites.IndicatorPopulation + sprites.IndicatorInactiveSuffix)
	s.indicatorProduction = s.sprites.GetIndex(sprites.IndicatorProduction)
	s.indicatorProductionInactive = s.sprites.GetIndex(sprites.IndicatorProduction + sprites.IndicatorInactiveSuffix)
	s.indicatorStorage = s.sprites.GetIndex(sprites.IndicatorStorage)
	s.indicatorStorageInactive = s.sprites.GetIndex(sprites.IndicatorStorage + sprites.IndicatorInactiveSuffix)

	fts := generic.NewResource[res.Fonts](world)
	fonts := fts.Get()
	s.font = fonts.Default
}

// UpdateUI the system
func (s *Terrain) UpdateUI(world *ecs.World) {
	sel := s.selection.Get()
	mouse := s.mouse.Get()
	ui := s.ui.Get()

	canvas := s.screen.Get()
	img := canvas.Image

	off := s.view.Offset()
	bounds := s.view.Bounds(canvas.Width, canvas.Height)

	img.Clear()
	img.Fill(s.sprites.Background)

	x, y := ebiten.CursorPosition()
	mx, my := s.view.ScreenToGlobal(x, y)
	cursor := s.view.GlobalToTile(mx, my)

	useMouse := mouse.IsInside && !ui.MouseInside(x, y)

	mapBounds := s.view.MapBounds(img.Bounds().Dx(), img.Bounds().Dy())
	mapBounds = mapBounds.Intersect(image.Rect(0, 0, s.terrain.Width(), s.terrain.Height()))

	showBuildable := sel.BuildType != terr.Air ||
		(s.landUse.Contains(cursor.X, cursor.Y) && terr.Properties[s.landUse.Get(cursor.X, cursor.Y)].BuildRadius > 0)

	for i := mapBounds.Min.X; i < mapBounds.Max.X; i++ {
		for j := mapBounds.Min.Y; j < mapBounds.Max.Y; j++ {
			point := s.view.TileToGlobal(i, j)
			if !point.In(bounds) {
				continue
			}

			height := 0
			t := s.terrain.Get(i, j)
			if t != terr.Air && t != terr.Buildable {
				tE := s.terrainE.Get(i, j)
				randTile := s.spriteMapper.Get(tE)
				height = s.drawSprite(img, s.terrain, s.landUse, i, j, t, &point, height, &off,
					randTile, terr.Properties[t].TerrainBelow, cursor.X, cursor.Y, sel.BuildType)

				if showBuildable {
					buildHere := s.buildable.Get(i, j) > 0
					buildMask, notBuildMask := s.buildable.NeighborsMask(i, j)
					if (buildHere && notBuildMask > 0) || (!buildHere && buildMask > 0) {
						if buildHere {
							_ = s.drawBorderSprite(img, s.borderInner, buildMask, &point, height, &off)
						} else {
							_ = s.drawBorderSprite(img, s.borderOuter, notBuildMask, &point, height, &off)
						}
					}
				}
			}

			lu := s.landUse.Get(i, j)
			if lu != terr.Air {
				luE := s.landUseE.Get(i, j)
				prod, randTile := s.landUseMapper.Get(luE)
				_ = s.drawSprite(img, s.terrain, s.landUse, i, j, lu, &point, height, &off,
					randTile, terr.Properties[lu].TerrainBelow, cursor.X, cursor.Y, sel.BuildType)
				if prod != nil &&
					(prod.Amount == 0 || prod.Stock >= terr.Properties[lu].Storage[prod.Resource]) {
					_ = s.drawSimpleSprite(img, s.warningMarker, &point, height, &off)
				}
			}

			if terr.Properties[lu].TerrainBits.Contains(terr.IsPath) {
				path := s.pathMapper.Get(s.landUseE.Get(i, j))
				offset := 0.1
				for _, h := range path.Haulers {
					haul, sp := s.haulerMapper.Get(h.Entity)
					s.drawHauler(img, sp.SpriteIndex, haul, height, offset, &off)
				}
			}

			if useMouse && cursor.X == i && cursor.Y == j {
				s.drawCursor(img, i, j, height, &point, &off, sel)
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
	sel *res.Selection) {

	ter := s.terrain.Get(x, y)
	lu := s.landUse.Get(x, y)
	luEntity := s.landUseE.Get(x, y)
	prop := terr.Properties[sel.BuildType]
	if prop.TerrainBits.Contains(terr.CanBuild) {
		canBuy := prop.TerrainBits.Contains(terr.CanBuy)

		canBuildHere := (prop.BuildOn.Contains(ter) || (sel.AllowRemove && ter != terr.Air && ter != sel.BuildType))
		isDestroy := false

		if prop.TerrainBits.Contains(terr.IsTerrain) {
			height = 0
			canBuildHere = canBuildHere && lu == terr.Air
			isDestroy = sel.AllowRemove && !prop.BuildOn.Contains(ter)
		} else {
			luNatural := !terr.Properties[lu].TerrainBits.Contains(terr.CanBuy)
			canBuildHere = canBuildHere &&
				(lu == terr.Air || (luNatural && canBuy)) &&
				(!prop.TerrainBits.Contains(terr.RequiresRange) || s.buildable.Get(x, y) > 0)
			isDestroy = lu != terr.Air && luNatural && canBuy
		}
		s.drawSprite(img, s.terrain, s.landUse, x, y, sel.BuildType, point, height, camOffset,
			&comp.RandomSprite{Rand: sel.RandSprite}, prop.TerrainBelow, x, y, terr.Air)

		cursor := s.cursorDenied
		if canBuildHere {
			if isDestroy || (sel.AllowRemove && !prop.BuildOn.Contains(ter)) {
				cursor = s.cursorDestroy
			} else {
				cursor = s.cursorOk
			}
		}
		s.drawCursorSprite(img, point, camOffset, cursor)
	} else if sel.BuildType == terr.Bulldoze {
		s.drawSprite(img, s.terrain, s.landUse, x, y, sel.BuildType, point, height, camOffset, nil, prop.TerrainBelow, x, y, terr.Air)
		if terr.Properties[lu].TerrainBits.Contains(terr.CanBuild) {
			s.drawCursorSprite(img, point, camOffset, s.cursorOk)
		} else {
			s.drawCursorSprite(img, point, camOffset, s.cursorDenied)
		}
	} else {
		s.drawCursorSprite(img, point, camOffset, s.cursorNeutral)
	}

	if sel.BuildType == terr.Air {
		s.drawBuildingMarker(img, lu, luEntity, point, camOffset)
	}
}

func (s *Terrain) drawBuildingMarker(img *ebiten.Image, lu terr.Terrain, e ecs.Entity, point, camOffset *image.Point) {
	if e.IsZero() {
		return
	}

	prop := &terr.Properties[lu]
	if prop.Production.MaxProduction > 0 {
		tp, prod := s.prodMapper.Get(e)
		bStock := terr.Properties[tp.Terrain].Storage[prod.Resource]

		h := s.sprites.Get(s.indicatorProduction).Bounds().Dy()

		s.drawIndicators(img, int(prod.Amount), int(prop.Production.MaxProduction),
			s.indicatorProduction, s.indicatorProductionInactive,
			point, camOffset, 0)
		s.drawIndicators(img, int(prod.Stock), int(bStock),
			s.indicatorStorage, s.indicatorStorageInactive,
			point, camOffset, -h)
	}
	if prop.PopulationSupport.MaxPopulation > 0 {
		pop := s.popMapper.Get(e)

		s.drawIndicators(img, int(pop.Pop), int(prop.PopulationSupport.MaxPopulation),
			s.indicatorPopulation, s.indicatorPopulationInactive,
			point, camOffset, 0)
	}
}

func (s *Terrain) drawIndicators(img *ebiten.Image,
	value, maxValue int, active, inactive int,
	point, camOffset *image.Point, yOffset int) {
	sp := s.sprites.Get(active)
	width := sp.Bounds().Dx()
	widthTotal := width * maxValue

	x := -widthTotal/2 + width/2
	for i := 0; i < value; i++ {
		pt := image.Pt(
			point.X+x,
			point.Y-yOffset-3*s.view.TileHeight,
		)
		s.drawSimpleSprite(img, active, &pt, 0, camOffset)
		x += width
	}
	for i := value; i < int(maxValue); i++ {
		pt := image.Pt(
			point.X+x,
			point.Y-yOffset-3*s.view.TileHeight,
		)
		s.drawSimpleSprite(img, inactive, &pt, 0, camOffset)
		x += width
	}
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
	below terr.Terrain,
	cursorX, cursorY int, cursorTerr terr.Terrain) int {

	idx := s.sprites.GetTerrainIndex(t)
	info := s.sprites.GetInfo(idx)

	if below != terr.Air {
		height = s.drawSprite(img, terrain, landUse,
			x, y, below, point, height,
			camOffset, randSprite, terr.Air, cursorX, cursorY, cursorTerr)
	}

	var sp *ebiten.Image
	if info.IsMultitile() {
		var neigh terr.Directions
		conn := terr.Properties[t].ConnectsTo
		cursorNear := terr.Properties[cursorTerr].TerrainBits.Contains(terr.CanBuild) &&
			math.AbsInt(x-cursorX) <= 1 && math.AbsInt(y-cursorY) <= 1
		if cursorNear {
			if terr.Properties[cursorTerr].TerrainBits.Contains(terr.IsTerrain) {
				neigh = terrain.NeighborsMaskMultiReplace(x, y, conn, cursorX, cursorY, cursorTerr) |
					landUse.NeighborsMaskMulti(x, y, conn)
			} else {
				neigh = terrain.NeighborsMaskMulti(x, y, conn) |
					landUse.NeighborsMaskMultiReplace(x, y, conn, cursorX, cursorY, cursorTerr)
			}
		} else {
			neigh = terrain.NeighborsMaskMulti(x, y, conn) |
				landUse.NeighborsMaskMulti(x, y, conn)
		}

		mIdx := s.sprites.GetMultiTileTerrainIndex(t, neigh, int(s.time.Tick), int(randSprite.GetRand()))

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

func (s *Terrain) drawBorderSprite(img *ebiten.Image, idx int,
	neigh terr.Directions, point *image.Point, height int,
	camOffset *image.Point) int {

	info := s.sprites.GetInfo(idx)

	var sp *ebiten.Image
	if info.IsMultitile() {
		mIdx := s.sprites.GetMultiTileIndex(idx, neigh, 0, 0)
		sp = s.sprites.GetSprite(mIdx)
	} else {
		sp = s.sprites.Get(idx)
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
	sp := s.sprites.GetRand(idx, int(s.time.Tick), 0)
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
