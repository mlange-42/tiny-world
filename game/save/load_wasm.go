//go:build js

package save

import (
	"encoding/json"
	"strings"
	"syscall/js"

	serde "github.com/mlange-42/ark-serde"
	"github.com/mlange-42/ark/ecs"
)

func loadWorld(world *ecs.World, folder, name string) error {
	_ = folder

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", saveGamePrefix+name)

	return serde.Deserialize([]byte(jsData.String()), world)
}

func loadSaveTime(folder, name string) (saveTime, error) {
	_ = folder

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", saveGamePrefix+name)

	helper := saveGameInfo{}
	err := json.Unmarshal([]byte(jsData.String()), &helper)
	if err != nil {
		return saveTime{}, err
	}
	return helper.Resources.SaveTime, nil
}

func loadAchievements(file string, completed *[]string) error {
	_ = file

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", achievementsKey)

	if jsData.IsNull() {
		return nil
	}

	return json.Unmarshal([]byte(jsData.String()), completed)
}

func listGames(folder string) ([]SaveGame, error) {
	_ = folder
	games := []SaveGame{}

	storage := js.Global().Get("localStorage")

	cnt := storage.Get("length").Int()
	for i := 0; i < cnt; i++ {
		key := storage.Call("key", i).String()
		if strings.HasPrefix(key, saveGamePrefix) {
			name := strings.TrimPrefix(key, saveGamePrefix)
			info, err := loadSaveTime(folder, name)
			if err != nil {
				return nil, err
			}

			games = append(games, SaveGame{
				Name: name,
				Time: info.Time,
			})
		}
	}

	return games, nil
}

func listMapsLocal(folder string) ([]MapLocation, error) {
	_ = folder
	maps := []MapLocation{}

	storage := js.Global().Get("localStorage")

	cnt := storage.Get("length").Int()
	for i := 0; i < cnt; i++ {
		key := storage.Call("key", i).String()
		if strings.HasPrefix(key, saveMapPrefix) {
			maps = append(maps, MapLocation{Name: strings.TrimPrefix(key, saveMapPrefix), IsEmbedded: false})
		}
	}

	return maps, nil
}

func loadMapLocal(folder string, name string) (string, error) {
	_ = folder
	storage := js.Global().Get("localStorage")
	mapData := storage.Call("getItem", saveMapPrefix+name)

	return mapData.String(), nil
}
