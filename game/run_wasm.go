//go:build js

package game

import (
	"embed"
	"syscall/js"

	serde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
)

func Run(data embed.FS) {
	gameData = data

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", "savegame")

	if jsData.IsNull() {
		run("", "paper")
	} else {
		if err := run("savegame", "paper"); err != nil {
			print(err.Error())
			run("", "paper")
		}
	}
}

func loadWorld(world *ecs.World, path string) error {
	_ = path

	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", path)

	return serde.Deserialize([]byte(jsData.String()), world)
}
