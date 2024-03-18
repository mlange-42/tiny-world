package save

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/maps"
	"github.com/mlange-42/tiny-world/game/terr"
)

func LoadWorld(world *ecs.World, folder, name string) error {
	_ = ecs.ComponentID[comp.Tile](world)
	_ = ecs.ComponentID[comp.Terrain](world)
	_ = ecs.ComponentID[comp.UpdateTick](world)
	_ = ecs.ComponentID[comp.Consumption](world)
	_ = ecs.ComponentID[comp.Production](world)
	_ = ecs.ComponentID[comp.Warehouse](world)
	_ = ecs.ComponentID[comp.BuildRadius](world)
	_ = ecs.ComponentID[comp.Path](world)
	_ = ecs.ComponentID[comp.Hauler](world)
	_ = ecs.ComponentID[comp.HaulerSprite](world)
	_ = ecs.ComponentID[comp.ProductionMarker](world)

	return loadWorld(world, folder, name)
}

func LoadAchievements(file string, completed *[]string) error {
	return loadAchievements(file, completed)
}

func ListSaveGames(folder string) ([]SaveGame, error) {
	games, err := listGames(folder)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(games, func(a, b SaveGame) int {
		return -a.Time.Compare(b.Time)
	})

	return games, nil
}

func LoadMap(f fs.FS, folder string, mapLoc MapLocation) (maps.Map, error) {
	mapStr, err := loadMap(f, folder, mapLoc)
	if err != nil {
		return maps.Map{}, err
	}

	return ParseMap(mapStr)
}

func ParseMap(mapStr string) (maps.Map, error) {
	helper := mapJs{}
	err := json.Unmarshal([]byte(mapStr), &helper)
	if err != nil {
		return maps.Map{}, nil
	}

	terrains := []rune{}
	for tStr, cnt := range helper.Terrains {
		rn := []rune(tStr)
		if len(rn) > 1 {
			return maps.Map{}, fmt.Errorf("symbols must be single runes. Got '%s'", tStr)
		}
		sym := rn[0]
		t, ok := terr.SymbolToTerrain[sym]
		if !ok {
			return maps.Map{}, fmt.Errorf("symbol not found: '%s'", tStr)
		}
		var ter terr.Terrain
		if t.LandUse != terr.Air {
			ter = t.LandUse
		} else {
			ter = t.Terrain
		}
		props := &terr.Properties[ter]
		if props.TerrainBits.Contains(terr.CanBuy) || !props.TerrainBits.Contains(terr.CanBuild) {
			return maps.Map{}, fmt.Errorf("terrain '%s' ('%s') is not a natural feature", props.Name, tStr)
		}
		for i := 0; i < cnt; i++ {
			terrains = append(terrains, sym)
		}
	}

	var result [][]rune
	for _, s := range helper.Map {
		if len(s) > 0 {
			runes := []rune(s)
			result = append(result, runes)
		}
	}

	return maps.Map{
		Data:                  result,
		Terrains:              terrains,
		InitialRandomTerrains: helper.InitialRandomTerrains,
		Center:                helper.Center,
		Achievements:          helper.Achievements,
		Description:           strings.Join(helper.Description, "\n"),
	}, nil
}

func LoadMapData(f fs.FS, folder string, mapLoc MapLocation) (MapInfo, error) {
	mapStr, err := loadMap(f, folder, mapLoc)
	if err != nil {
		return MapInfo{}, err
	}
	helper := mapInfoJs{}
	err = json.Unmarshal([]byte(mapStr), &helper)
	if err != nil {
		return MapInfo{}, nil
	}

	return MapInfo{Achievements: helper.Achievements, Description: strings.Join(helper.Description, "\n")}, nil
}

func ListMaps(f fs.FS, folder string) ([]MapLocation, error) {
	lst, err := listMapsEmbed(f, folder)
	if err != nil {
		return nil, err
	}
	lst2, err := listMapsLocal(folder)
	if err != nil {
		return nil, err
	}
	return append(lst, lst2...), nil
}

func loadMap(f fs.FS, folder string, mapLoc MapLocation) (string, error) {
	if mapLoc.IsEmbedded {
		mapData, err := fs.ReadFile(f, path.Join("data", folder, mapLoc.Name)+".json")
		if err != nil {
			return "", err
		}

		return string(mapData), nil
	}

	return loadMapLocal(folder, mapLoc.Name)
}

func listMapsEmbed(f fs.FS, folder string) ([]MapLocation, error) {
	maps := []MapLocation{}

	files, err := fs.ReadDir(f, path.Join("data", folder))
	if err != nil {
		return nil, nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == ".json" {
			base := strings.TrimSuffix(file.Name(), ".json")
			maps = append(maps, MapLocation{Name: base, IsEmbedded: true})
		}
	}
	return maps, nil
}
