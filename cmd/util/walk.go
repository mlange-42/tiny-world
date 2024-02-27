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

type SpriteSheet struct {
	Directory string
	Width     int
	Height    int
}

type Directory struct {
	Dir         string
	HasJson     bool
	Files       []File
	Directories []Directory
}

type File struct {
	Name   string
	IsJson bool
}

func Walk(base string, tileSet string, f func(sheet SpriteSheet, dir Directory)) error {
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
			f(sheet, info)
		}
	}

	return nil
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

func newTileSheet(dir string) (SpriteSheet, error) {
	parts := strings.Split(dir, "_")
	if len(parts) != 2 {
		return SpriteSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	size := strings.Split(parts[1], "x")
	if len(size) != 2 {
		return SpriteSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	w, err := strconv.Atoi(size[0])
	if err != nil {
		return SpriteSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	h, err := strconv.Atoi(size[1])
	if err != nil {
		return SpriteSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}

	return SpriteSheet{Directory: dir, Width: w, Height: h}, nil
}
