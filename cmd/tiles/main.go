package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
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
	tmath "github.com/mlange-42/tiny-world/math"
)

const maxWidth = 512

const suffixMultiTile = "_multitile"

const outFolder = "assets/sprites"
const inFolder = "artwork/sprites"

func main() {
	dirs := extractFiles()

	for _, dir := range dirs {
		processDirectory(dir)
	}
}

func processDirectory(info dirInfo) {
	fmt.Printf("Processing %s (%d images)\n", info.Directory, len(info.Files))

	if len(info.Files) == 0 {
		return
	}

	mask := isoMask(info.Width, info.Height)

	infos := []util.SpriteInfo{}
	images := []image.Image{}

	index := 0
	for _, file := range info.Files {
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

		spriteInfo := util.SpriteInfo{
			Name:  baseName,
			Index: index,
		}
		jsonFile := strings.Replace(file, filepath.Ext(file), "", 1) + ".json"
		if content, err := os.ReadFile(jsonFile); err == nil {
			if err := json.Unmarshal(content, &spriteInfo); err != nil {
				log.Fatal("error decoding JSON: ", err)
			}
		}

		if strings.HasSuffix(baseName, suffixMultiTile) {
			if sprite.Bounds().Dx() != info.Width*4 || sprite.Bounds().Dy() != info.Height*4 {
				log.Fatalf("unexpected tile size in %s: got %dx%d", file, sprite.Bounds().Dx(), sprite.Bounds().Dy())
			}
			tiles := spiltMultiTile(sprite, mask, info.Width, info.Height)
			for i, tile := range tiles {
				name := strings.Replace(baseName, suffixMultiTile, "", 1)
				if i > 0 {
					name = fmt.Sprintf("%s_%d", name, i)
				}
				spriteInfo := util.SpriteInfo{
					Name:      name,
					Index:     index,
					MultiTile: i == 0,
				}

				images = append(images, tile)
				infos = append(infos, spriteInfo)
				index++
			}
		} else {
			if sprite.Bounds().Dx() != info.Width || sprite.Bounds().Dy() != info.Height {
				log.Fatalf("unexpected tile size in %s: got %dx%d", file, sprite.Bounds().Dx(), sprite.Bounds().Dy())
			}

			spriteInfo.MultiTile = false

			images = append(images, sprite)
			infos = append(infos, spriteInfo)

			index++
		}
	}

	perRow := maxWidth / info.Width
	numRows := int(math.Ceil(float64(len(images)) / float64(perRow)))

	sheetWidth := perRow * info.Width
	sheetHeight := numRows * info.Height

	img := image.NewRGBA(image.Rect(0, 0, sheetWidth, sheetHeight))
	for i, sprite := range images {
		row, col := i/perRow, i%perRow
		draw.Draw(img,
			image.Rect(col*info.Width, row*info.Height, col*info.Width+info.Width, row*info.Height+info.Height),
			sprite, image.Point{}, draw.Src)
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

func isoMask(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	fill := color.RGBA{255, 255, 255, 255}
	mask := color.RGBA{0, 0, 0, 0}

	draw.Draw(img, img.Bounds(), &image.Uniform{fill}, image.Point{}, draw.Src)

	midX := width / 2
	midY := height / 2

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			dx := (midX - x)
			dy := (midY - y)

			if x >= midX {
				dx = x + 1 - midX
			}
			if y >= midY {
				dy = y + 1 - midY
			}

			dist := dx + 2*dy
			if dist > height+1 {
				img.SetRGBA(x, y, mask)
			}
		}
	}
	return img
}

func spiltMultiTile(sprite image.Image, mask *image.RGBA, width, height int) []*image.RGBA {
	result := []*image.RGBA{}

	dx, dy := width/2, height/2
	doubleSize := 8

	for row := 0; row < doubleSize-1; row++ {
		perRow := tmath.MinInt(row, doubleSize-2-row) + 1
		halfOffsets := (doubleSize - 2*perRow) / 2
		xOffset := halfOffsets * dx
		yOffset := row * dy
		for col := 0; col < perRow; col++ {
			img := image.NewRGBA(image.Rect(0, 0, width, height))
			draw.DrawMask(img, img.Bounds(), sprite, image.Point{xOffset + col*width, yOffset}, mask, image.Point{}, draw.Src)
			result = append(result, img)
		}
	}

	return result
}
