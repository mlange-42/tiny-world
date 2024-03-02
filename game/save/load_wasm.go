//go:build js

package save

import (
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

func listSaveGames(folder string) ([]string, error) {
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
