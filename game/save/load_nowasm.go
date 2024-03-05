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

func listSaveGames(folder string) ([]string, error) {
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
		if ext == ".json" || ext == ".JSON" {
			base := strings.TrimSuffix(file.Name(), ".json")
			base = strings.TrimSuffix(base, ".JSON")
			games = append(games, base)
		}
	}
	return games, nil
}
