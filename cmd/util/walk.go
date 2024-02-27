package util

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const multiTileSuffix = "_multitile"

func WalkSheets(base string, tileSet string, f func(sheet RawSpriteSheet) error) error {
	basePath := path.Join(base, tileSet)
	sheets, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, sheetDir := range sheets {
		if !sheetDir.IsDir() {
			continue
		}
		sheet, err := newTileSheet(sheetDir.Name())
		if err != nil {
			log.Printf("SKIP: %s\n", err.Error())
		}
		if err := f(sheet); err != nil {
			return err
		}
	}
	return nil
}

func WalkDirs(base string, tileSet string, sheet RawSpriteSheet, f func(sheet RawSpriteSheet, dir Directory) error) error {
	sheetPath := path.Join(base, tileSet, sheet.Directory)
	dirs, err := os.ReadDir(sheetPath)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		info := walkDir(sheetPath, dir.Name(), true)
		if err := f(sheet, info); err != nil {
			return err
		}
	}
	return nil
}

func Walk(base string, tileSet string, f func(sheet RawSpriteSheet, dir Directory) error) error {
	return WalkSheets(base, tileSet, func(sheet RawSpriteSheet) error {
		return WalkDirs(base, tileSet, sheet, f)
	})
}

func walkDir(base string, dir string, recursive bool) Directory {
	p := path.Join(base, dir)
	files, err := os.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	info := Directory{Dir: dir}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == ".json" || ext == ".JSON" {
			info.HasJson = true
			break
		}
	}
	for _, file := range files {
		if file.IsDir() {
			if recursive && strings.HasSuffix(file.Name(), multiTileSuffix) {
				subInfo := walkDir(p, file.Name(), false)
				info.Directories = append(info.Directories, subInfo)
			}
			continue
		}
		ext := filepath.Ext(file.Name())
		if (info.HasJson && (ext == ".json" || ext == ".JSON")) ||
			(!info.HasJson && (ext == ".png" || ext == ".PNG")) {
			info.Files = append(info.Files, File{file.Name(), info.HasJson})
		}
	}

	return info
}

func newTileSheet(dir string) (RawSpriteSheet, error) {
	parts := strings.Split(dir, "_")
	if len(parts) != 2 {
		return RawSpriteSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	size := strings.Split(parts[1], "x")
	if len(size) != 2 {
		return RawSpriteSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	w, err := strconv.Atoi(size[0])
	if err != nil {
		return RawSpriteSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	h, err := strconv.Atoi(size[1])
	if err != nil {
		return RawSpriteSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}

	return RawSpriteSheet{Directory: dir, Width: w, Height: h}, nil
}
