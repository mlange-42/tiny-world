//go:build !js

package save

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"

	serde "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
)

func loadWorld(world *ecs.World, folder, name string) error {
	jsData, err := os.ReadFile(path.Join(folder, name) + ".json")
	if err != nil {
		return err
	}

	return serde.Deserialize(jsData, world)
}

func loadSaveTime(folder, name string) (saveTime, error) {
	jsData, err := os.ReadFile(path.Join(folder, name) + ".json")
	if err != nil {
		return saveTime{}, err
	}
	helper := saveGameInfo{}
	err = json.Unmarshal(jsData, &helper)
	if err != nil {
		return saveTime{}, err
	}
	return helper.Resources.SaveTime, nil
}

func loadAchievements(file string, completed *[]string) error {
	jsData, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsData, completed)
}

func listGames(folder string) ([]SaveGame, error) {
	games := []SaveGame{}

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
			info, err := loadSaveTime(folder, base)
			if err != nil {
				return nil, err
			}
			games = append(games, SaveGame{
				Name: base,
				Time: info.Time,
			})
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
		if ext == ".json" {
			base := strings.TrimSuffix(file.Name(), ".json")
			maps = append(maps, MapLocation{Name: base, IsEmbedded: false})
		}
	}
	return maps, nil
}

func loadMapLocal(folder string, name string) (string, error) {
	mapData, err := os.ReadFile(path.Join(folder, name) + ".json")
	if err != nil {
		return "", err
	}

	return string(mapData), nil
}
