package save

import (
	"fmt"
	"regexp"
	"strings"

	as "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

const mapDescriptionDelimiter = "----"

func SaveWorld(folder, name string, world *ecs.World, skip []generic.Comp) error {
	js, err := as.Serialize(world,
		as.Opts.SkipResources(
			skip...,
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

	terrains := map[rune]int{}

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
		if cnt, ok := terrains[sym]; ok {
			terrains[sym] = cnt + 1
		} else {
			terrains[sym] = 1
		}
	}

	symbols := []string{}
	for sym, cnt := range terrains {
		symbols = append(symbols, fmt.Sprintf("%d%s", cnt, string(sym)))
	}
	b.WriteString(strings.Join(symbols, " "))
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("%d\n", rules.InitialRandomTerrains))

	// Space for required achievements
	b.WriteString("\n")

	// Delimiter for map description
	b.WriteString(mapDescriptionDelimiter + "\n")

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
