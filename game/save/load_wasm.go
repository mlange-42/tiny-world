//go:build js

package save

import (
	"strings"
	"syscall/js"

	serde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
)

const (
	fileTypeJson  fileType = saveGamePrefix
	fileTypeAscii fileType = saveMapPrefix
)

func loadWorld(world *ecs.World, folder, name string) error {
	_ = folder

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", saveGamePrefix+name)

	return serde.Deserialize([]byte(jsData.String()), world)
}

func listFiles(folder string, ft fileType) ([]string, error) {
	_ = folder
	games := []string{}

	storage := js.Global().Get("localStorage")

	cnt := storage.Get("length").Int()
	for i := 0; i < cnt; i++ {
		key := storage.Call("key", i).String()
		if strings.HasPrefix(key, string(ft)) {
			games = append(games, strings.TrimPrefix(key, string(ft)))
		}
	}

	return games, nil
}

func loadMap(folder, name string) (string, error) {
	_ = folder

	storage := js.Global().Get("localStorage")
	mapData := storage.Call("getItem", saveMapPrefix+name)

	return mapData.String(), nil
}
