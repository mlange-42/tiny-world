//go:build js

package save

import (
	"syscall/js"

	serde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
)

func loadWorld(world *ecs.World, path string) error {
	_ = path

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", path)

	return serde.Deserialize([]byte(jsData.String()), world)
}

func listSaveGames(folder string) ([]string, error) {
	_ = folder
	games := []string{}

	storage := js.Global().Get("localStorage")

	cnt := storage.Get("length").Int()
	for i := 0; i < cnt; i++ {
		key := storage.Call("key", i).String()
		games = append(games, key)
	}

	return games, nil
}
