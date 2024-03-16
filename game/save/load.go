package save

import (
	"fmt"
	"image"
	"io/fs"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/maps"
)

type LoadType uint8

type MapLocation struct {
	Name       string
	IsEmbedded bool
}

const (
	LoadTypeNone LoadType = iota
	LoadTypeGame
	LoadTypeMap
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

func ListSaveGames(folder string) ([]string, error) {
	return listGames(folder)
}

func LoadMap(f fs.FS, folder string, mapLoc MapLocation) (maps.Map, error) {
	mapStr, err := loadMap(f, folder, mapLoc)
	if err != nil {
		return maps.Map{}, err
	}

	return ParseMap(mapStr)
}

func ParseMap(mapStr string) (maps.Map, error) {
	var result [][]rune
	lines := strings.Split(strings.ReplaceAll(mapStr, "\r\n", "\n"), "\n")

	terrainParts := strings.Split(lines[0], " ")
	terrains := []rune{}
	for _, p := range terrainParts {
		rn := []rune(p)
		sym := rn[len(rn)-1]
		cnt, err := strconv.Atoi(string(rn[:len(rn)-1]))
		if err != nil {
			panic(fmt.Sprintf("can't convert to integer in map symbols: `%s`", string(rn[:len(rn)-1])))
		}
		for i := 0; i < cnt; i++ {
			terrains = append(terrains, sym)
		}
	}

	randTerr, err := strconv.Atoi(lines[1])
	if err != nil {
		panic(fmt.Sprintf("can't convert to integer in map initial random terrains: `%s`", lines[1]))
	}

	ach := strings.Split(lines[2], " ")
	achievements := []string{}
	for _, a := range ach {
		if a != "" {
			achievements = append(achievements, a)
		}
	}

	sizeLine := lines[3]
	parts := strings.Split(sizeLine, " ")
	cx, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(fmt.Sprintf("can't convert to integer: `%s`", parts[0]))
	}
	cy, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(fmt.Sprintf("can't convert to integer: `%s`", parts[1]))
	}

	lines = lines[4:]

	for _, s := range lines {
		if len(s) > 0 {
			runes := []rune(s)
			result = append(result, runes)
		}
	}

	return maps.Map{
		Data:                  result,
		Terrains:              terrains,
		InitialRandomTerrains: randTerr,
		Center:                image.Pt(cx, cy),
		Achievements:          achievements,
	}, nil
}

func LoadMapAchievements(f fs.FS, folder string, mapLoc MapLocation) ([]string, error) {
	mapStr, err := loadMap(f, folder, mapLoc)
	if err != nil {
		return nil, err
	}

	lines := strings.SplitN(strings.ReplaceAll(mapStr, "\r\n", "\n"), "\n", 4)

	ach := strings.Split(lines[2], " ")
	achievements := []string{}
	for _, a := range ach {
		if a != "" {
			achievements = append(achievements, a)
		}
	}

	return achievements, nil
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
		mapData, err := fs.ReadFile(f, path.Join("data", folder, mapLoc.Name)+".asc")
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
		if ext == ".asc" {
			base := strings.TrimSuffix(file.Name(), ".asc")
			maps = append(maps, MapLocation{Name: base, IsEmbedded: true})
		}
	}
	return maps, nil
}
