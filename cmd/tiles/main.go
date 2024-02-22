package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mlange-42/tiny-world/game/util"
)

const maxWidth = 512

const outFolder = "assets/sprites"
const inFolder = "artwork/sprites"

func main() {
	dirs := extractFiles()

	for _, dir := range dirs {
		processDirectory(dir)
	}
}

func processDirectory(info dirInfo) {
	count := len(info.Files)

	fmt.Printf("Processing %s (%d images)\n", info.Directory, count)

	if count == 0 {
		return
	}

	perRow := maxWidth / info.Width
	numRows := int(math.Ceil(float64(count) / float64(perRow)))

	width := perRow * info.Width
	height := numRows * info.Height

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	infos := []util.SpriteInfo{}
	for i, file := range info.Files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("error reading file %s: %s", file, err.Error())
		}
		defer f.Close()
		sprite, _, err := image.Decode(f)
		if err != nil {
			log.Fatalf("error decoding file %s: %s", file, err.Error())
		}

		baseName := strings.Replace(filepath.Base(file), filepath.Ext(file), "", 1)

		row, col := i/perRow, i%perRow

		draw.Draw(img,
			image.Rect(col*info.Width, row*info.Height, col*info.Width+info.Width, row*info.Height+info.Height),
			sprite, image.Point{}, draw.Src)

		infos = append(infos, util.SpriteInfo{
			Name:  baseName,
			Index: i,
		})
	}

	outFile := path.Join(outFolder, fmt.Sprintf("%s.png", info.Directory))
	outFileJson := path.Join(outFolder, fmt.Sprintf("%s.json", info.Directory))

	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("error creating file %s: %s", outFile, err.Error())
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("error encoding image: %s", err.Error())
	}

	spriteSheet := util.SpriteSheet{
		SpriteWidth:  info.Width,
		SpriteHeight: info.Height,
		Sprites:      infos,
	}
	js, err := json.MarshalIndent(spriteSheet, "", " ")
	if err != nil {
		log.Fatalf("error encoding json: %s", err.Error())
	}
	if err := os.WriteFile(outFileJson, js, 0666); err != nil {
		log.Fatalf("error writing JSON file: %s", err.Error())
	}

}

type dirInfo struct {
	Directory string
	Width     int
	Height    int
	Files     []string
}

func extractFiles() []dirInfo {
	entries, err := os.ReadDir(inFolder)
	if err != nil {
		log.Fatal(err)
	}

	dirs := []dirInfo{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		info, err := newDirInfo(e.Name())
		if err != nil {
			log.Printf("SKIP: %s\n", err.Error())
		}
		dirs = append(dirs, info)
	}

	for i := range dirs {
		dir := &dirs[i]
		p := path.Join(inFolder, dir.Directory)
		entries, err := os.ReadDir(p)
		if err != nil {
			log.Fatal(err)
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			p := path.Join(inFolder, dir.Directory, e.Name())
			files, err := os.ReadDir(p)
			if err != nil {
				log.Fatal(err)
			}

			for _, e := range files {
				if e.IsDir() {
					continue
				}
				ext := filepath.Ext(e.Name())
				if ext != ".png" && ext != ".PNG" {
					continue
				}
				dir.Files = append(dir.Files, path.Join(p, e.Name()))
			}
		}
	}

	return dirs
}

func newDirInfo(dir string) (dirInfo, error) {
	parts := strings.Split(dir, "_")
	if len(parts) != 2 {
		return dirInfo{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	size := strings.Split(parts[1], "x")
	if len(size) != 2 {
		return dirInfo{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	w, err := strconv.Atoi(size[0])
	if err != nil {
		return dirInfo{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}
	h, err := strconv.Atoi(size[1])
	if err != nil {
		return dirInfo{}, fmt.Errorf("directory does not match expected pattern: %s", dir)
	}

	return dirInfo{Directory: dir, Width: w, Height: h}, nil
}
