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
	return g.Contains(x+dx, y+dy) && g.Get(x+dx, y+dy) == tp
}

func (g *TerrainGrid) isNeighborMask(x, y, dx, dy int, tp terr.Terrains) bool {
	return g.Contains(x+dx, y+dy) && tp.Contains(g.Get(x+dx, y+dy))
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
