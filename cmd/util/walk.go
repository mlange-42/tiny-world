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

type TileSheet struct {
	Directory string
	Width     int
	Height    int
}

type Directory struct {
	Dir         string
	HasJson     bool
	Files       []File
	Directories []string
}

type File struct {
	Name   string
	IsJson bool
}

func Walk(base string, tileSet string, f func(sheet TileSheet, dir Directory)) error {
	basePath := path.Join(base, tileSet)
	sheets, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}
	fmt.Println(sheets)

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

			dirPath := path.Join(sheetPath, dir.Name())
			files, err := os.ReadDir(dirPath)
			if err != nil {
				log.Fatal(err)
			}

			info := Directory{Dir: dir.Name()}
			for _, file := range files {
				if file.IsDir() {
					if strings.HasSuffix(file.Name(), multiTileSuffix) {
						info.Directories = append(info.Directories, file.Name())
					}
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
					continue
				}
				ext := filepath.Ext(file.Name())
				if (info.HasJson && (ext == ".json" || ext == ".JSON")) ||
					(!info.HasJson && (ext == ".png" || ext == ".PNG")) {
					info.Files = append(info.Files, File{file.Name(), info.HasJson})
				}
			}

			f(sheet, info)
		}
	}

	return nil
}

func newTileSheet(dir string) (TileSheet, error) {
	parts := strings.Split(dir, "_")
	if len(parts) != 2 {
		return TileSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	size := strings.Split(parts[1], "x")
	if len(size) != 2 {
		return TileSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	w, err := strconv.Atoi(size[0])
	if err != nil {
		return TileSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	h, err := strconv.Atoi(size[1])
	if err != nil {
		return TileSheet{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}

	return TileSheet{Directory: dir, Width: w, Height: h}, nil
}
