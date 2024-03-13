package save

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mlange-42/arche-model/model"
	"github.com/mlange-42/arche-model/resource"
	as "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
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
			generic.T[res.Achievements](),
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

func SaveAchievements(file string, completed []string) error {
	return saveAchievements(file, completed)
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

func SaveMap(folder, name string, world *ecs.World) error {
	b := strings.Builder{}

	rules := ecs.GetResource[res.Rules](world)
	bounds := ecs.GetResource[res.WorldBounds](world)
	terrain := ecs.GetResource[res.Terrain](world)
	landUse := ecs.GetResource[res.LandUse](world)

	for _, t := range rules.RandomTerrains {
		var tp terr.TerrainPair
		if terr.Properties[t].TerrainBits.Contains(terr.IsTerrain) {
			tp.Terrain = t
		} else {
			tp.LandUse = t
		}
		sym, ok := terr.TerrainToSymbol[tp]
		if !ok {
			return fmt.Errorf("symbol not found for %s/%s", terr.Properties[tp.Terrain].Name, terr.Properties[tp.LandUse].Name)
		}
		b.WriteRune(sym)
	}
	b.WriteString("\n")

	cx, cy := terrain.Width()/2, terrain.Height()/2
	b.WriteString(fmt.Sprintf("%d %d\n", cx-bounds.Min.X, cy-bounds.Min.Y))

	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			ter := terrain.Get(x, y)
			if !terr.Properties[ter].TerrainBits.Contains(terr.CanBuild) {
				ter = terr.Air
			}
			t := terr.TerrainPair{Terrain: ter, LandUse: landUse.Get(x, y)}
			sym, ok := terr.TerrainToSymbol[t]
			if !ok {
				return fmt.Errorf("symbol not found for %s/%s", terr.Properties[t.Terrain].Name, terr.Properties[t.LandUse].Name)
			}
			b.WriteRune(sym)
		}
		b.WriteString("\n")
	}

	return saveMapToFile(folder, name, b.String())
}
