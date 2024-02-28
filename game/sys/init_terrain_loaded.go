package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// InitTerrainLoaded system.
type InitTerrainLoaded struct {
}

// Initialize the system
func (s *InitTerrainLoaded) Initialize(world *ecs.World) {
	terrain := ecs.GetResource[res.Terrain](world)
	terrainE := ecs.GetResource[res.TerrainEntities](world)
	landUse := ecs.GetResource[res.LandUse](world)
	landUseE := ecs.GetResource[res.LandUseEntities](world)

	filter := generic.NewFilter2[comp.Tile, comp.Terrain]()
	query := filter.Query(world)
	for query.Next() {
		tile, ter := query.Get()
		if terr.Properties[ter.Terrain].IsTerrain {
			terrain.Set(tile.X, tile.Y, ter.Terrain)
			terrainE.Set(tile.X, tile.Y, query.Entity())
		} else {
			landUse.Set(tile.X, tile.Y, ter.Terrain)
			landUseE.Set(tile.X, tile.Y, query.Entity())
		}
	}
}

// Update the system
func (s *InitTerrainLoaded) Update(world *ecs.World) {}

// Finalize the system
func (s *InitTerrainLoaded) Finalize(world *ecs.World) {}
