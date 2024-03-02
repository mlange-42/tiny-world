//go:build js

package save

import "syscall/js"

func saveToFile(path string, jsData []byte) error {
	_ = path

	data := js.ValueOf(string(jsData))
	storage := js.Global().Get("localStorage")
	storage.Set("savegame", data)

	return nil
}
