//go:build js

package save

import "syscall/js"

func saveToFile(folder, name string, jsData []byte) error {
	_ = folder

	data := js.ValueOf(string(jsData))
	storage := js.Global().Get("localStorage")
	storage.Call("setItem", name, data)

	return nil
}
