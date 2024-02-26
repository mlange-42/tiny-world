package nav

import (
	"container/heap"
	"image"

	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/math"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

type Score struct {
	Tile  comp.Tile
	Score int
}

type AStar struct {
	landUse *res.LandUse
}

func NewAStar(landUse *res.LandUse) AStar {
	return AStar{
		landUse: landUse,
	}
}

func (a *AStar) FindPath(start, target comp.Tile) ([]comp.Tile, bool) {
	initialDist := math.AbsInt(start.X-target.X) + math.AbsInt(start.Y-target.Y)

	open := NewPriorityQueue()
	heap.Init(&open)
	heap.Push(&open, Score{start, initialDist})

	cameFrom := map[comp.Tile]comp.Tile{}
	gScore := map[comp.Tile]int{}
	gScore[start] = 0

	for open.Len() > 0 {
		current := heap.Pop(&open).(Score)
		if current.Tile == target {
			return reconstruct(cameFrom, current.Tile), true
		}
		luOld := a.landUse.Get(current.Tile.X, current.Tile.Y)
		if current.Tile != start {
			if luOld != terr.Path {
				gScore[current.Tile] = 0
				continue
			}
		}

		for dir := terr.Direction(0); dir < terr.EndDirection; dir++ {
			dx, dy := dir.Deltas()
			xx, yy := current.Tile.X+dx, current.Tile.Y+dy
			if !a.landUse.Contains(xx, yy) {
				continue
			}
			lu := a.landUse.Get(xx, yy)
			if !(lu == terr.Path || (luOld == terr.Path && terr.Buildings.Contains(lu))) {
				continue
			}

			other := comp.Tile{Point: image.Pt(xx, yy)}
			heur := gScore[current.Tile] + math.AbsInt(dx) + math.AbsInt(dy)

			otherScore := 999999999
			if sc, ok := gScore[other]; ok {
				otherScore = sc
			}
			if heur < otherScore {
				cameFrom[other] = current.Tile

				fSc := heur + math.AbsInt(other.X-target.X) + math.AbsInt(other.Y-target.Y)
				gScore[other] = heur
				if open.Contains(other) {
					open.Update(other, fSc)
				} else {
					heap.Push(&open, Score{other, fSc})
				}
			}
		}
	}

	return nil, false
}

func reconstruct(cameFrom map[comp.Tile]comp.Tile, current comp.Tile) []comp.Tile {
	path := []comp.Tile{current}
	for {
		if v, ok := cameFrom[current]; ok {
			current = v
			path = append(path, v)
		} else {
			break
		}
	}
	return path
}

// An TileHeap is a min-heap of ints.
type TileHeap struct {
	score []Score
	tiles map[comp.Tile]bool
}

func NewTileMap() TileHeap {
	return TileHeap{
		tiles: map[comp.Tile]bool{},
	}
}

func (h TileHeap) Len() int           { return len(h.score) }
func (h TileHeap) Less(i, j int) bool { return h.score[i].Score < h.score[j].Score }
func (h TileHeap) Swap(i, j int)      { h.score[i], h.score[j] = h.score[j], h.score[i] }

func (h *TileHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	score := x.(Score)
	h.score = append(h.score, score)
	h.tiles[score.Tile] = true
}

func (h *TileHeap) Pop() any {
	old := h.score
	n := len(old)
	x := old[n-1]
	h.score = old[0 : n-1]
	delete(h.tiles, x.Tile)
	return x
}

func (h *TileHeap) Contains(tile comp.Tile) bool {
	_, ok := h.tiles[tile]
	return ok
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue struct {
	score []Score
	tiles map[comp.Tile]int
}

func NewPriorityQueue() PriorityQueue {
	return PriorityQueue{
		tiles: map[comp.Tile]int{},
	}
}

func (pq PriorityQueue) Len() int { return len(pq.score) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq.score[i].Score < pq.score[j].Score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.score[i], pq.score[j] = pq.score[j], pq.score[i]
	pq.tiles[pq.score[i].Tile] = i
	pq.tiles[pq.score[j].Tile] = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(pq.score)
	item := x.(Score)
	pq.tiles[item.Tile] = n
	pq.score = append(pq.score, item)
}

func (pq *PriorityQueue) Pop() any {
	old := pq.score
	n := len(old)
	item := old[n-1]
	delete(pq.tiles, item.Tile)
	pq.score = old[0 : n-1]
	return item
}

// Update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(tile comp.Tile, priority int) {
	idx := pq.tiles[tile]
	pq.score[idx].Score = priority
	heap.Fix(pq, idx)
}

func (h *PriorityQueue) Contains(tile comp.Tile) bool {
	_, ok := h.tiles[tile]
	return ok
}
