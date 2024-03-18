//go:build js

package save

import (
	"encoding/json"
	"syscall/js"
)

// TODO change to pattern to match saveMapPrefix, on the next save-game breaking change.
const saveGamePrefix = "tiny-world-save-"
const saveMapPrefix = "mlange-42/tiny-world/maps/"
const achievementsKey = "mlange-42/tiny-world/achievements"

func saveToFile(folder, name string, jsData []byte) error {
	_ = folder

	data := js.ValueOf(string(jsData))
	storage := js.Global().Get("localStorage")
	storage.Call("setItem", saveGamePrefix+name, data)

	return nil
}

func saveAchievements(file string, completed []string) error {
	_ = file

	jsData, err := json.MarshalIndent(completed, "", " ")
	if err != nil {
		return err
	}

	data := js.ValueOf(string(jsData))
	storage := js.Global().Get("localStorage")
	storage.Call("setItem", achievementsKey, data)

	return nil
}

func deleteGame(folder, name string) error {
	_ = folder

	storage := js.Global().Get("localStorage")
	storage.Delete(saveGamePrefix + name)
	return nil
}

func saveMapToFile(folder, name string, mapData string) error {
	_ = folder

	data := js.ValueOf(mapData)
	storage := js.Global().Get("localStorage")
	storage.Call("setItem", saveMapPrefix+name, data)

	return nil
}
