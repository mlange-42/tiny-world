package save

import (
	"encoding/json"
	"fmt"
	"image"
	"regexp"
	"strings"

	as "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

func SaveWorld(folder, name string, world *ecs.World, skip []ecs.Comp) error {
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
	rules := ecs.GetResource[res.Rules](world)
	bounds := ecs.GetResource[res.WorldBounds](world)
	terrain := ecs.GetResource[res.Terrain](world)
	landUse := ecs.GetResource[res.LandUse](world)

	terrains := map[string]int{}

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
		symStr := string(sym)
		if cnt, ok := terrains[symStr]; ok {
			terrains[symStr] = cnt + 1
		} else {
			terrains[symStr] = 1
		}
	}

	center := image.Pt(terrain.Width()/2-bounds.Min.X, terrain.Height()/2-bounds.Min.Y)

	rows := make([]string, bounds.Dy()+1)
	b := strings.Builder{}
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
		rows[y-bounds.Min.Y] = b.String()
		b.Reset()
	}

	mapJs := mapJs{
		Terrains:              terrains,
		Center:                center,
		InitialRandomTerrains: rules.InitialRandomTerrains,
		Achievements:          []string{},
		Description:           []string{},
		Map:                   rows,
	}
	jsData, err := json.MarshalIndent(mapJs, "", "  ")
	if err != nil {
		return err
	}

	return saveMapToFile(folder, name, string(jsData))
}
