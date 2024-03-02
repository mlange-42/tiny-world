//go:build js

package save

import "syscall/js"

const saveGamePrefix = "tiny-world-save-"

func saveToFile(folder, name string, jsData []byte) error {
	_ = folder

	data := js.ValueOf(string(jsData))
	storage := js.Global().Get("localStorage")
	storage.Call("setItem", saveGamePrefix+name, data)

	return nil
}
