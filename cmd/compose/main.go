package main

import (
	"fmt"
	"image"
	"path"

	"github.com/mlange-42/tiny-world/cmd/util"
)

const (
	basePath = "artwork"
	tileSet  = "sprites"
)

func main() {
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

	err := util.WalkSheets(basePath, tileSet, func(sheet util.SpriteSheet) error {

		clear(p.Indices)
		clear(p.Images)
		clear(p.Infos)
		p.Images = p.Images[:0]
		p.Infos = p.Infos[:0]

		err := util.WalkDirs(basePath, tileSet, sheet,
			func(sheet util.SpriteSheet, dir util.Directory) error {
				if dir.HasJson {
					p.processDirectoryJson(sheet, dir)
				} else {
					p.processDirectoryNoJson(sheet, dir)
				}
				return nil
			},
		)
		if err != nil {
			return err
		}
		fmt.Println(p.Infos)

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (p *proc) processDirectoryJson(sheet util.SpriteSheet, dir util.Directory) {
	base := path.Join(basePath, tileSet, sheet.Directory, dir.Dir)
	sprites := []util.RawSprite{}
	files := []string{}

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

			files = append(files, js.File...)
			for _, f := range js.Multitile {
				files = append(files, f...)
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
			Id:        sprite.Id,
			Height:    sprite.Height,
			YOffset:   sprite.YOffset,
			Index:     make([]int, len(sprite.File)),
			Multitile: make([][]int, len(sprite.Multitile)),
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

func (p *proc) processDirectoryNoJson(sheet util.SpriteSheet, dir util.Directory) {
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
			Id:    file.Name,
			Index: []int{index},
		})
	}
}
