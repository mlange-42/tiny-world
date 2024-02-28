package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
	"os"
	"path"
	"strings"

	"github.com/mlange-42/tiny-world/cmd/util"
	"github.com/spf13/cobra"
)

const (
	basePath = "artwork/sprites"
	outPath  = "assets/sprites"

	maxSheetWidth = 512
)

func main() {
	if err := command().Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(tileSet string) {
	p := proc{}
	p.Process(basePath, tileSet)
}

type proc struct {
	Names   map[string]bool
	Images  []image.Image
	Infos   []util.Sprite
	Indices map[string]int
}

func (p *proc) Process(basePath, tileSet string) {
	p.Names = map[string]bool{}
	p.Indices = map[string]int{}

	tilesetJs := util.TileSet{}
	if err := util.FromJson(path.Join(basePath, tileSet, "tileset.json"), &tilesetJs); err != nil {
		panic(err)
	}
	if err := util.ToJson(path.Join(outPath, tileSet, "tileset.json"), &tilesetJs); err != nil {
		panic(err)
	}

	err := util.WalkSheets(basePath, tileSet, func(sheet util.RawSpriteSheet) error {

		clear(p.Indices)
		clear(p.Images)
		clear(p.Infos)
		p.Images = p.Images[:0]
		p.Infos = p.Infos[:0]

		if err := util.WalkDirs(basePath, tileSet, sheet,
			func(sheet util.RawSpriteSheet, dir util.Directory) error {
				if dir.HasJson {
					p.processDirectoryJson(tileSet, sheet, dir)
				} else {
					p.processDirectoryNoJson(tileSet, sheet, dir)
				}
				return nil
			},
		); err != nil {
			return err
		}

		if err := p.writeSheet(outPath, tileSet, sheet, maxSheetWidth); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (p *proc) writeSheet(outPath, tileSet string, sheet util.RawSpriteSheet, maxWidth int) error {
	outPathBase := path.Join(outPath, tileSet, sheet.Directory)

	outSheet := util.SpriteSheet{
		SpriteWidth:  sheet.Width,
		SpriteHeight: sheet.Height,
		Sprites:      p.Infos,
		TotalSprites: len(p.Images),
	}

	err := util.ToJson(outPathBase+".json", &outSheet)
	if err != nil {
		return err
	}

	perRow := maxWidth / sheet.Width
	numRows := int(math.Ceil(float64(len(p.Images)) / float64(perRow)))

	sheetWidth := perRow * sheet.Width
	sheetHeight := numRows * sheet.Height

	w, h := sheet.Width, sheet.Height

	img := image.NewRGBA(image.Rect(0, 0, sheetWidth, sheetHeight))
	for i, sprite := range p.Images {
		row, col := i/perRow, i%perRow
		draw.Draw(img,
			image.Rect(col*w, row*h, col*w+w, row*h+h),
			sprite, image.Point{}, draw.Over)
	}

	return util.WriteImage(outPathBase+".png", img)
}

func (p *proc) processDirectoryJson(tileSet string, sheet util.RawSpriteSheet, dir util.Directory) {
	base := path.Join(basePath, tileSet, sheet.Directory, dir.Dir)
	sprites := []util.RawSprite{}
	files := []string{}
	filesDone := map[string]bool{}

	for _, file := range dir.Files {
		filePath := path.Join(base, file.Name)

		jsArr := []util.RawSprite{}
		err := util.FromJson(filePath, &jsArr)
		if err != nil {
			panic(err)
		}
		for _, js := range jsArr {
			if js.Id == "" {
				panic(fmt.Errorf("missing ID in '%s'", filePath))
			}
			if len(js.File) == 0 {
				js.File = []string{js.Id}
			}
			if _, ok := p.Names[js.Id]; ok {
				panic(fmt.Errorf("duplicate sprite ID '%s' in %s", js.Id, filePath))
			}
			p.Names[js.Id] = true
			p.Indices[js.Id] = len(p.Images)

			for _, f := range js.File {
				if _, ok := filesDone[f]; ok {
					continue
				}
				filesDone[f] = true
				files = append(files, f)
			}
			for _, fs := range js.Multitile {
				for _, f := range fs {
					if _, ok := filesDone[f]; ok {
						continue
					}
					filesDone[f] = true
					files = append(files, f)
				}
			}

			sprites = append(sprites, js)
		}
	}

	indices := map[string]int{}
	zeroIndex := len(p.Images)

	for i, file := range files {
		img, err := util.ReadImage(path.Join(base, file) + ".png")
		if err != nil {
			panic(err)
		}
		p.Images = append(p.Images, img)
		indices[file] = zeroIndex + i
	}

	for _, sprite := range sprites {
		sp := util.Sprite{
			Id:         sprite.Id,
			Height:     sprite.Height,
			YOffset:    sprite.YOffset,
			AnimFrames: sprite.AnimFrames,
			AnimSpeed:  sprite.AnimSpeed,
			Index:      make([]int, len(sprite.File)),
			Multitile:  make([][]int, len(sprite.Multitile)),
		}

		for i, id := range sprite.File {
			sp.Index[i] = indices[id]
		}
		for i, ids := range sprite.Multitile {
			sp.Multitile[i] = make([]int, len(ids))
			for j, id := range ids {
				sp.Multitile[i][j] = indices[id]
			}
		}
		p.Infos = append(p.Infos, sp)
	}
}

func (p *proc) processDirectoryNoJson(tileSet string, sheet util.RawSpriteSheet, dir util.Directory) {
	base := path.Join(basePath, tileSet, sheet.Directory, dir.Dir)

	for _, file := range dir.Files {
		if _, ok := p.Names[file.Name]; ok {
			panic(fmt.Sprintf("duplicate sprite ID '%s'", file.Name))
		}
		index := len(p.Images)
		p.Names[file.Name] = true
		p.Indices[file.Name] = index

		img, err := util.ReadImage(path.Join(base, file.Name))
		if err != nil {
			panic(err)
		}
		p.Images = append(p.Images, img)
		p.Infos = append(p.Infos, util.Sprite{
			Id:    strings.ReplaceAll(file.Name, ".png", ""),
			Index: []int{index},
		})
	}
}

func command() *cobra.Command {
	var tileSet string
	root := &cobra.Command{
		Use:           "go run ./cmd/compose",
		Short:         "Compose sprite sheets",
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
