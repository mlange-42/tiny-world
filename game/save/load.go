package save

import (
	"strings"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
)

type fileType string

type LoadType uint8

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

func ListSaveGames(folder string) ([]string, error) {
	return listFiles(folder, fileTypeJson)
}

func LoadMap(folder, name string) ([][]rune, error) {
	mapStr, err := loadMap(folder, name)
	if err != nil {
		return nil, err
	}

	var result [][]rune
	lines := strings.Split(strings.ReplaceAll(mapStr, "\r\n", "\n"), "\n")
	for _, s := range lines {
		if len(s) > 0 {
			result = append(result, []rune(s))
		}
	}

	return result, nil
}

func ListMaps(folder string) ([]string, error) {
	return listFiles(folder, fileTypeAscii)
}
