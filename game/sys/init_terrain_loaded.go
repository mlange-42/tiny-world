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
	rules := ecs.GetResource[res.Rules](world)
	terrain := ecs.GetResource[res.Terrain](world)
	terrainE := ecs.GetResource[res.TerrainEntities](world)
	landUse := ecs.GetResource[res.LandUse](world)
	landUseE := ecs.GetResource[res.LandUseEntities](world)
	fac := ecs.GetResource[res.EntityFactory](world)

	filter := generic.NewFilter2[comp.Tile, comp.Terrain]()
	query := filter.Query(world)
	for query.Next() {
		tile, ter := query.Get()
		if terr.Properties[ter.Terrain].TerrainBits.Contains(terr.IsTerrain) {
			terrain.Set(tile.X, tile.Y, ter.Terrain)
			terrainE.Set(tile.X, tile.Y, query.Entity())
		} else {
			landUse.Set(tile.X, tile.Y, ter.Terrain)
			landUseE.Set(tile.X, tile.Y, query.Entity())
		}
	}

	x, y := terrain.Width()/2, terrain.Height()/2
	fac.SetBuildable(x, y, rules.InitialBuildRadius, true)

	radFilter := generic.NewFilter2[comp.Tile, comp.BuildRadius]()
	radQuery := radFilter.Query(world)
	for radQuery.Next() {
		tile, rad := radQuery.Get()
		fac.SetBuildable(tile.X, tile.Y, int(rad.Radius), true)
	}
}

// Update the system
func (s *InitTerrainLoaded) Update(world *ecs.World) {}

// Finalize the system
func (s *InitTerrainLoaded) Finalize(world *ecs.World) {}
