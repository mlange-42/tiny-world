package res

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Terrain resource
type Terrain struct {
	TerrainGrid
}

func NewTerrain(w, h int) Terrain {
	return Terrain{
		TerrainGrid: TerrainGrid{NewGrid[terr.Terrain](w, h)},
	}
}

// LandUse resource
type LandUse struct {
	TerrainGrid
}

func NewLandUse(w, h int) LandUse {
	return LandUse{
		TerrainGrid: TerrainGrid{NewGrid[terr.Terrain](w, h)},
	}
}

// TerrainEntities resource
type TerrainEntities struct {
	Grid[ecs.Entity]
}

// LandUseEntities resource
type LandUseEntities struct {
	Grid[ecs.Entity]
}

// Buildable resource
type Buildable struct {
	Grid[uint16]
}

func NewBuildable(w, h int) Buildable {
	return Buildable{
		Grid: NewGrid[uint16](w, h),
	}
}

func (b *Buildable) NeighborsMask(x, y int) (terr.Directions, terr.Directions) {
	dirs := terr.Directions(0)
	notDirs := terr.Directions(0)
	if b.isNeighbor(x, y, 0, -1) {
		dirs.Set(terr.N)
	} else {
		notDirs.Set(terr.N)
	}
	if b.isNeighbor(x, y, 1, 0) {
		dirs.Set(terr.E)
	} else {
		notDirs.Set(terr.E)
	}
	if b.isNeighbor(x, y, 0, 1) {
		dirs.Set(terr.S)
	} else {
		notDirs.Set(terr.S)
	}
	if b.isNeighbor(x, y, -1, 0) {
		dirs.Set(terr.W)
	} else {
		notDirs.Set(terr.W)
	}
	return dirs, notDirs
}

func (b *Buildable) isNeighbor(x, y, dx, dy int) bool {
	return b.Contains(x+dx, y+dy) && b.Get(x+dx, y+dy) != 0
}

type TerrainGrid struct {
	Grid[terr.Terrain]
}

func (g *TerrainGrid) NeighborsMaskMulti(x, y int, tp terr.Terrains) terr.Directions {
	dirs := terr.Directions(0)
	if g.isNeighborMask(x, y, 0, -1, tp) {
		dirs.Set(terr.N)
	}
	if g.isNeighborMask(x, y, 1, 0, tp) {
		dirs.Set(terr.E)
	}
	if g.isNeighborMask(x, y, 0, 1, tp) {
		dirs.Set(terr.S)
	}
	if g.isNeighborMask(x, y, -1, 0, tp) {
		dirs.Set(terr.W)
	}
	return dirs
}

func (g *TerrainGrid) NeighborsMaskMultiReplace(x, y int, tp terr.Terrains, rx, ry int, rt terr.Terrain) terr.Directions {
	dirs := terr.Directions(0)
	if g.isNeighborMaskReplace(x, y, 0, -1, tp, rx, ry, rt) {
		dirs.Set(terr.N)
	}
	if g.isNeighborMaskReplace(x, y, 1, 0, tp, rx, ry, rt) {
		dirs.Set(terr.E)
	}
	if g.isNeighborMaskReplace(x, y, 0, 1, tp, rx, ry, rt) {
		dirs.Set(terr.S)
	}
	if g.isNeighborMaskReplace(x, y, -1, 0, tp, rx, ry, rt) {
		dirs.Set(terr.W)
	}
	return dirs
}

func (g *TerrainGrid) NeighborsMask(x, y int, tp terr.Terrain) terr.Directions {
	dirs := terr.Directions(0)
	if g.isNeighbor(x, y, 0, -1, tp) {
		dirs.Set(terr.N)
	}
	if g.isNeighbor(x, y, 1, 0, tp) {
		dirs.Set(terr.E)
	}
	if g.isNeighbor(x, y, 0, 1, tp) {
		dirs.Set(terr.S)
	}
	if g.isNeighbor(x, y, -1, 0, tp) {
		dirs.Set(terr.W)
	}
	return dirs
}

func (g *TerrainGrid) isNeighbor(x, y, dx, dy int, tp terr.Terrain) bool {
	xx, yy := x+dx, y+dy
	return g.Contains(xx, yy) && g.Get(xx, yy) == tp
}

func (g *TerrainGrid) isNeighborMask(x, y, dx, dy int, tp terr.Terrains) bool {
	xx, yy := x+dx, y+dy
	return g.Contains(xx, yy) && tp.Contains(g.Get(xx, yy))
}

func (g *TerrainGrid) isNeighborMaskReplace(x, y, dx, dy int, tp terr.Terrains, rx, ry int, rt terr.Terrain) bool {
	xx, yy := x+dx, y+dy
	if xx == rx && yy == ry {
		return tp.Contains(rt)
	}
	return g.Contains(xx, yy) && tp.Contains(g.Get(xx, yy))
}

func (g *TerrainGrid) CountNeighbors4(x, y int, tp terr.Terrain) int {
	cnt := 0
	if g.isNeighbor(x, y, 0, -1, tp) {
		cnt++
	}
	if g.isNeighbor(x, y, 1, 0, tp) {
		cnt++
	}
	if g.isNeighbor(x, y, 0, 1, tp) {
		cnt++
	}
	if g.isNeighbor(x, y, -1, 0, tp) {
		cnt++
	}
	return cnt
}

func (g *TerrainGrid) CountNeighborsMask4(x, y int, tp terr.Terrains) int {
	cnt := 0
	if g.isNeighborMask(x, y, 0, -1, tp) {
		cnt++
	}
	if g.isNeighborMask(x, y, 1, 0, tp) {
		cnt++
	}
	if g.isNeighborMask(x, y, 0, 1, tp) {
		cnt++
	}
	if g.isNeighborMask(x, y, -1, 0, tp) {
		cnt++
	}
	return cnt
}

func (g *TerrainGrid) CountNeighbors8(x, y int, tp terr.Terrain) int {
	cnt := g.CountNeighbors4(x, y, tp)
	if g.isNeighbor(x, y, 1, -1, tp) {
		cnt++
	}
	if g.isNeighbor(x, y, 1, 1, tp) {
		cnt++
	}
	if g.isNeighbor(x, y, -1, 1, tp) {
		cnt++
	}
	if g.isNeighbor(x, y, -1, -1, tp) {
		cnt++
	}
	return cnt
}

func (g *TerrainGrid) CountNeighborsMask8(x, y int, tp terr.Terrains) int {
	cnt := g.CountNeighborsMask4(x, y, tp)
	if g.isNeighborMask(x, y, 1, -1, tp) {
		cnt++
	}
	if g.isNeighborMask(x, y, 1, 1, tp) {
		cnt++
	}
	if g.isNeighborMask(x, y, -1, 1, tp) {
		cnt++
	}
	if g.isNeighborMask(x, y, -1, -1, tp) {
		cnt++
	}
	return cnt
}
