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

func LoadMap(f fs.FS, folder, name string) ([][]rune, image.Point, error) {
	mapStr, err := loadMap(f, folder, name)
	if err != nil {
		return nil, image.Point{}, err
	}

	var result [][]rune
	lines := strings.Split(strings.ReplaceAll(mapStr, "\r\n", "\n"), "\n")

	sizeLine := lines[0]
	parts := strings.Split(sizeLine, " ")
	cx, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(fmt.Sprintf("can't convert to integer: `%s`", parts[0]))
	}
	cy, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(fmt.Sprintf("can't convert to integer: `%s`", parts[1]))
	}

	lines = lines[1:]

	for _, s := range lines {
		if len(s) > 0 {
			runes := []rune(s)
			result = append(result, runes)
		}
	}

	return result, image.Pt(cx, cy), nil
}

func ListMaps(f fs.FS, folder string) ([]string, error) {
	return listFilesFS(f, folder, fileTypeAscii)
}

func loadMap(f fs.FS, folder, name string) (string, error) {
	mapData, err := fs.ReadFile(f, path.Join(folder, name)+".asc")
	if err != nil {
		return "", err
	}

	return string(mapData), nil
}

func listFilesFS(f fs.FS, folder string, ft fileType) ([]string, error) {
	games := []string{}

	files, err := fs.ReadDir(f, folder)
	if err != nil {
		return nil, nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == string(ft) {
			base := strings.TrimSuffix(file.Name(), string(ft))
			games = append(games, base)
		}
	}
	return games, nil
}
