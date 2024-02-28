package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"path"

	"github.com/mlange-42/tiny-world/cmd/util"
	gMath "github.com/mlange-42/tiny-world/game/math"
	"github.com/spf13/cobra"
)

const (
	basePath = "artwork/sprites"
)

type multitileJson struct {
	Id     string `json:"id"`
	File   string `json:"file"`
	Base   string `json:"base"`
	Below  string `json:"below"`
	Height int    `json:"height"`
}

var multiTileOrder = [16]int{
	4,
	5, 6,
	1, 7, 14,
	0, 3, 15, 12,
	2, 11, 13,
	10, 9,
	8,
}

func main() {
	if err := command().Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

}

func run(tileSet string) {
	if err := util.Walk(basePath, tileSet, func(sheet util.RawSpriteSheet, dir util.Directory) error {
		for _, subDir := range dir.Directories {
			if subDir.HasJson {
				processDirectoryJson(tileSet, sheet, dir, subDir)
			} else {
				processDirectoryNoJson(tileSet, sheet, dir, subDir)
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func processDirectoryJson(tileSet string, sheet util.RawSpriteSheet, dir, subDir util.Directory) {
	base := path.Join(basePath, tileSet, sheet.Directory, dir.Dir)
	sub := path.Join(base, subDir.Dir)
	for _, file := range subDir.Files {
		filePath := path.Join(sub, file.Name)
		js := multitileJson{}
		err := util.FromJson(filePath, &js)
		if err != nil {
			log.Fatal("error reading multitile JSON: ", err)
		}
		if js.Id == "" {
			log.Fatal("error reading multitile JSON: missing ID in ", filePath)
		}
		if js.File == "" {
			js.File = js.Id
		}
		processFile(base, subDir.Dir, js)
	}
}

func processDirectoryNoJson(tileSet string, sheet util.RawSpriteSheet, dir, subDir util.Directory) {
	base := path.Join(basePath, tileSet, sheet.Directory, dir.Dir)
	for _, file := range subDir.Files {
		js := multitileJson{
			Id:   file.Name,
			File: file.Name,
		}
		processFile(base, subDir.Dir, js)
	}
}

func processFile(base, src string, js multitileJson) {
	filePath := path.Join(base, src, js.File) + ".png"

	img, err := util.ReadImage(filePath)
	if err != nil {
		log.Fatal(err)
	}

	w, h := img.Bounds().Dx()/4, img.Bounds().Dy()/4
	mask := isoMask(w, h)
	images := spiltMultiTile(img, mask, w, h)

	if js.Base != "" {
		baseFilePath := path.Join(base, src, js.Base) + ".png"
		baseImage, err := util.ReadImage(baseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		for i := range images {
			newImage := image.NewRGBA(baseImage.Bounds())
			draw.Draw(newImage, newImage.Rect, baseImage, image.Point{}, draw.Over)
			draw.Draw(newImage, newImage.Rect, images[i], image.Point{}, draw.Over)
			images[i] = newImage
		}
	}

	outJs := util.RawSprite{
		Id:     js.Id,
		File:   []string{js.File},
		Height: js.Height,
		Below:  js.Below,
	}

	for i := range images {
		var name string
		if i == 0 {
			name = js.File
		} else {
			name = fmt.Sprintf("%s_%d", js.File, i)
		}
		outPath := path.Join(base, fmt.Sprintf("%s.png", name))
		err := util.WriteImage(outPath, images[i])
		if err != nil {
			log.Fatal(err)
		}

		outJs.Multitile = append(outJs.Multitile, []string{name})
	}

	err = util.ToJson(path.Join(base, fmt.Sprintf("%s.json", js.File)), &[]util.RawSprite{outJs})
	if err != nil {
		log.Fatal(err)
	}
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
	result := make([]*image.RGBA, len(multiTileOrder))

	dx, dy := width/2, height/2
	doubleSize := 8

	index := 0
	for row := 0; row < doubleSize-1; row++ {
		perRow := gMath.MinInt(row, doubleSize-2-row) + 1
		halfOffsets := (doubleSize - 2*perRow) / 2
		xOffset := halfOffsets * dx
		yOffset := row * dy
		for col := 0; col < perRow; col++ {
			img := image.NewRGBA(image.Rect(0, 0, width, height))
			draw.DrawMask(img, img.Bounds(), sprite, image.Point{xOffset + col*width, yOffset}, mask, image.Point{}, draw.Src)
			result[multiTileOrder[index]] = img
			index++
		}
	}

	return result
}

func command() *cobra.Command {
	var tileSet string
	root := &cobra.Command{
		Use:           "go run ./cmd/slice",
		Short:         "Slice multi-tiles",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if tileSet == "" {
				_ = cmd.Help()
				log.Fatal("please provide a tileset!")
			}
			run(tileSet)
		},
	}
	root.Flags().StringVarP(&tileSet, "tileset", "t", "", "Tileset to process.")

	return root
}
