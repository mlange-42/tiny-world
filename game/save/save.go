package save

import (
	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/resource"
	as "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

func SaveWorld(path string, world *ecs.World) error {
	js, err := as.Serialize(world,
		as.Opts.SkipResources(
			generic.T[res.Fonts](),
			generic.T[res.EbitenImage](),
			generic.T[res.EntityFactory](),
			generic.T[res.Sprites](),
			generic.T[res.Terrain](),
			generic.T[res.TerrainEntities](),
			generic.T[res.LandUse](),
			generic.T[res.LandUseEntities](),
			generic.T[resource.Termination](),
			generic.T[model.Systems](),
		),
	)
	if err != nil {
		return err
	}

	return saveToFile(path, js)
}
