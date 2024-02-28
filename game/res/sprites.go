package res

import (
	"encoding/json"
	"fmt"
	"image"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mlange-42/tiny-world/cmd/util"
	"github.com/mlange-42/tiny-world/game/terr"
)

const nameUnknown = "unknown"

type Sprites struct {
	atlas       []*ebiten.Image
	sprites     []*ebiten.Image
	infos       []util.Sprite
	indices     map[string]int
	terrIndices [terr.EndTerrain]int
	idxUnknown  int
}

func NewSprites(fSys fs.FS, dir string) Sprites {
	sheets, err := fs.ReadDir(fSys, dir)
	if err != nil {
		log.Fatal("error reading sprites", err)
	}

	atlas := []*ebiten.Image{}
	sprites := []*ebiten.Image{}
	infos := []util.Sprite{}
	indices := map[string]int{}

	infoIndex := 0
	imageIndex := 0
	for _, sheetFile := range sheets {
		if sheetFile.IsDir() {
			continue
		}
		ext := filepath.Ext(sheetFile.Name())
		if ext != ".json" && ext != ".JSON" {
			continue
		}
		baseName := strings.Replace(sheetFile.Name(), ext, "", 1)
		pngPath := path.Join(dir, fmt.Sprintf("%s.png", baseName))

		sheet := util.SpriteSheet{}
		content, err := fs.ReadFile(fSys, path.Join(dir, sheetFile.Name()))
		if err != nil {
			log.Fatal("error loading JSON file: ", err)
		}
		if err := json.Unmarshal(content, &sheet); err != nil {
			log.Fatal("error decoding JSON: ", err)
		}

		img, _, err := ebitenutil.NewImageFromFileSystem(fSys, pngPath)
		if err != nil {
			log.Fatal("error reading image: ", err)
		}
		atlas = append(atlas, img)

		fmt.Printf("%s -- %d sprites, %d images\n", baseName, len(sheet.Sprites), sheet.TotalSprites)
		for _, inf := range sheet.Sprites {
			if _, ok := indices[inf.Id]; ok {
				log.Fatalf("duplicate sprite name: %s", inf.Id)
			}
			indices[inf.Id] = infoIndex

			for i := range inf.Index {
				inf.Index[i] += imageIndex
			}
			for i := range inf.Multitile {
				for j := range inf.Multitile[i] {
					inf.Multitile[i][j] += imageIndex
				}
			}
			infos = append(infos, inf)
			infoIndex++
		}

		w, h := sheet.SpriteWidth, sheet.SpriteHeight
		cols, _ := img.Bounds().Dx()/w, img.Bounds().Dy()/h

		for i := 0; i < sheet.TotalSprites; i++ {
			row := i / cols
			col := i % cols
			sprites = append(sprites, img.SubImage(image.Rect(col*w, row*h, col*w+w, row*h+h)).(*ebiten.Image))
		}

		imageIndex += sheet.TotalSprites
	}

	terrIndices := [terr.EndTerrain]int{}
	for i := terr.Terrain(0); i < terr.EndTerrain; i++ {
		if idx, ok := indices[terr.Properties[i].Name]; ok {
			terrIndices[i] = idx
		} else {
			terrIndices[i] = indices[nameUnknown]
		}
	}

	return Sprites{
		atlas:       atlas,
		sprites:     sprites,
		infos:       infos,
		indices:     indices,
		idxUnknown:  indices[nameUnknown],
		terrIndices: terrIndices,
	}
}

func (s *Sprites) Get(idx int) (*ebiten.Image, *util.Sprite) {
	inf := &s.infos[idx]
	return s.sprites[inf.Index[0]], inf
}

func (s *Sprites) GetRand(idx int, rand int) (*ebiten.Image, *util.Sprite) {
	inf := &s.infos[idx]
	return s.sprites[inf.Index[rand%len(inf.Index)]], inf
}

func (s *Sprites) GetSprite(idx int) *ebiten.Image {
	return s.sprites[idx]
}

func (s *Sprites) GetIndex(name string) int {
	if idx, ok := s.indices[name]; ok {
		return idx
	}
	return s.idxUnknown
}

func (s *Sprites) GetTerrainIndex(t terr.Terrain) int {
	return s.terrIndices[t]
}

func (s *Sprites) GetMultiTileIndex(t terr.Terrain, dirs terr.Directions) int {
	idx := s.terrIndices[t]
	if s.infos[idx].IsMultitile() {
		return s.infos[idx].Multitile[dirs][0]
	}
	return idx
}

func (s *Sprites) GetMultiTileIndexRand(t terr.Terrain, dirs terr.Directions, rand int) int {
	idx := s.terrIndices[t]
	if s.infos[idx].IsMultitile() {
		sprites := s.infos[idx].Multitile[dirs]
		return sprites[rand%len(sprites)]
	}
	return idx
}
