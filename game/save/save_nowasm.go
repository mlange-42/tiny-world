//go:build !js

package save

import (
	"os"
	"path/filepath"
)

func saveToFile(path string, jsData []byte) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(jsData)

	return nil
}
