package res

import "github.com/mlange-42/tiny-world/game/terr"

// Grid data structure
type Grid[T comparable] struct {
	data   []T
	width  int
	height int
}

// NewGrid returns a new Grid.
func NewGrid[T comparable](width, height int) Grid[T] {
	return Grid[T]{
		data:   make([]T, width*height),
		width:  width,
		height: height,
	}
}

// Get a value from the Grid.
func (g *Grid[T]) Get(x, y int) T {
	idx := x*g.height + y
	return g.data[idx]
}

// Get a pointer to a value from the Grid.
func (g *Grid[T]) GetPointer(x, y int) *T {
	idx := x*g.height + y
	return &g.data[idx]
}

// Set a value in the grid.
func (g *Grid[T]) Set(x, y int, value T) {
	idx := x*g.height + y
	g.data[idx] = value
}

// Fill the grid with a value.
func (g *Grid[T]) Fill(value T) {
	for i := range g.data {
		g.data[i] = value
	}
}

// Width of the Grid.
func (g *Grid[T]) Width() int {
	return g.width
}

// Height of the Grid.
func (g *Grid[T]) Height() int {
	return g.height
}

// Contains returns whether the grid contains the given cell.
func (g *Grid[T]) Contains(x, y int) bool {
	return x >= 0 && y >= 0 && x < g.width && y < g.height
}

// Clamp coordinates to be inside the grid.
func (g *Grid[T]) Clamp(x, y int) (int, int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= g.width {
		x = g.width - 1
	}
	if y >= g.height {
		y = g.height - 1
	}
	return x, y
}

func (g *Grid[T]) NeighborsMask(x, y int, tp T) terr.Directions {
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

func (g *Grid[T]) isNeighbor(x, y, dx, dy int, tp T) bool {
	return g.Contains(x+dx, y+dy) && g.Get(x+dx, y+dy) == tp
}
