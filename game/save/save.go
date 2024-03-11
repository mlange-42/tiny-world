package save

import (
	"regexp"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/resource"
	as "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

func SaveWorld(folder, name string, world *ecs.World) error {
	js, err := as.Serialize(world,
		as.Opts.SkipResources(
			generic.T[res.Fonts](),
			generic.T[res.Screen](),
			generic.T[res.EntityFactory](),
			generic.T[res.Sprites](),
			generic.T[res.Terrain](),
			generic.T[res.TerrainEntities](),
			generic.T[res.LandUse](),
			generic.T[res.LandUseEntities](),
			generic.T[res.Buildable](),
			generic.T[res.SaveEvent](),
			generic.T[res.UpdateInterval](),
			generic.T[res.GameSpeed](),
			generic.T[res.Selection](),
			generic.T[res.Mouse](),
			generic.T[res.View](),
			generic.T[resource.Termination](),
			generic.T[resource.Rand](),
			generic.T[model.Systems](),
		),
	)
	if err != nil {
		return err
	}

	return saveToFile(folder, name, js)
}

func IsValidName(name string) bool {
	re := `^[a-zA-Z0-9][a-zA-Z0-9 \-_]*$`
	matched, err := regexp.Match(re, []byte(name))
	if err != nil {
		panic(err)
	}
	return matched
}

func DeleteGame(folder, name string) error {
	return deleteGame(folder, name)
}
