package res

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Terrain resource
type Terrain struct {
	Grid[terr.Terrain]
}

func (t *Terrain) Set(x, y int, value terr.Terrain) {
	t.Grid.Set(x, y, value)
	t.setNeighbor(x, y, -1, 0)
	t.setNeighbor(x, y, 1, 0)
	t.setNeighbor(x, y, 0, -1)
	t.setNeighbor(x, y, 0, 1)
}

func (t *Terrain) setNeighbor(x, y, dx, dy int) {
	if t.Contains(x+dx, y+dy) && t.Get(x+dx, y+dy) == terr.Air {
		t.Grid.Set(x+dx, y+dy, terr.Buildable)
	}
}

// LandUse resource
type LandUse struct {
	Grid[terr.Terrain]
}

// LandUseEntities resource
type LandUseEntities struct {
	Grid[ecs.Entity]
}
