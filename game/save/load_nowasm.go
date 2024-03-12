//go:build !js

package save

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	serde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
)

func loadWorld(world *ecs.World, folder, name string) error {
	jsData, err := os.ReadFile(path.Join(folder, name) + ".json")
	if err != nil {
		return err
	}

	return serde.Deserialize(jsData, world)
}

func listGames(folder string) ([]string, error) {
	games := []string{}

	files, err := os.ReadDir(folder)
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
			games = append(games, base)
		}
	}
	return games, nil
}

func listMapsLocal(folder string) ([]MapLocation, error) {
	maps := []MapLocation{}

	files, err := os.ReadDir(folder)
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
			maps = append(maps, MapLocation{Name: base, IsEmbedded: false})
		}
	}
	return maps, nil
}

func loadMapLocal(folder string, name string) (string, error) {
	mapData, err := os.ReadFile(path.Join(folder, name) + ".asc")
	if err != nil {
		return "", err
	}

	return string(mapData), nil
}
