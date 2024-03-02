//go:build js

package game

import "syscall/js"

func run(g *Game) {
	storage := js.Global().Get("localStorage")
	jsData := storage.Call("getItem", "savegame")

	if jsData.IsNull() {
		runGame(g, "", "paper")
	} else {
		if err := runGame(g, "savegame", "paper"); err != nil {
			print(err.Error())
			runGame(g, "", "paper")
		}
	}
}
