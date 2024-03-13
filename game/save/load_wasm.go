//go:build js

package save

import (
	"encoding/json"
	"strings"
	"syscall/js"

	serde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
)

func loadWorld(world *ecs.World, folder, name string) error {
	_ = folder

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", saveGamePrefix+name)

	return serde.Deserialize([]byte(jsData.String()), world)
}

func loadAchievements(file string, completed *[]string) error {
	_ = file

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", achievementsKey)

	return json.Unmarshal([]byte(jsData.String()), completed)
}

func listGames(folder string) ([]string, error) {
	_ = folder
	games := []string{}

	storage := js.Global().Get("localStorage")

	cnt := storage.Get("length").Int()
	for i := 0; i < cnt; i++ {
		key := storage.Call("key", i).String()
		if strings.HasPrefix(key, saveGamePrefix) {
			games = append(games, strings.TrimPrefix(key, saveGamePrefix))
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
