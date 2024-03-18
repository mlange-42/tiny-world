//go:build !js

package save

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

func saveToFile(folder, name string, jsData []byte) error {
	file := path.Join(folder, name) + ".json"
	dir := filepath.Dir(file)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(jsData)

	return nil
}

func saveAchievements(file string, completed []string) error {
	jsData, err := json.MarshalIndent(completed, "", " ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(file)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(jsData)

	return nil
}

func deleteGame(folder, name string) error {
	file := path.Join(folder, name) + ".json"
	return os.Remove(file)
}

func saveMapToFile(folder, name string, mapData string) error {
	file := path.Join(folder, name) + ".json"
	dir := filepath.Dir(file)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(mapData)

	return nil
}
