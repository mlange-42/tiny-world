//go:build js

package game

import (
	"embed"
	"syscall/js"
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
