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
	"github.com/mlange-42/tiny-world/game/terr"
	"github.com/mlange-42/tiny-world/game/util"
)

const nameUnknown = "unknown"

type Sprites struct {
	atlas       []*ebiten.Image
	sprites     []*ebiten.Image
	infos       []Sprite
	indices     map[string]int
	terrIndices [terr.EndTerrain]int
	idxUnknown  int
}

type Sprite struct {
	Height    int
	YOffset   int
	MultiTile bool
}

func NewSprites(fSys fs.FS, dir string) Sprites {
	entries, err := fs.ReadDir(fSys, dir)
	if err != nil {
		log.Fatal("error reading sprites", err)
	}

	atlas := []*ebiten.Image{}
	sprites := []*ebiten.Image{}
	infos := []Sprite{}
	indices := map[string]int{}

	index := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext != ".json" && ext != ".JSON" {
			continue
		}
		baseName := strings.Replace(e.Name(), ext, "", 1)
		pngPath := path.Join(dir, fmt.Sprintf("%s.png", baseName))

		sheet := util.SpriteSheet{}
		content, err := fs.ReadFile(fSys, path.Join(dir, e.Name()))
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

		w, h := sheet.SpriteWidth, sheet.SpriteHeight
		cols, _ := img.Bounds().Dx()/w, img.Bounds().Dy()/h

		fmt.Println(e.Name())
		for i, inf := range sheet.Sprites {
			if _, ok := indices[inf.Name]; ok {
				log.Fatalf("duplicate sprite name: %s", inf.Name)
			}
			indices[inf.Name] = index

			row := i / cols
			col := i % cols
			sprites = append(sprites, img.SubImage(image.Rect(col*w, row*h, col*w+w, row*h+h)).(*ebiten.Image))

			infos = append(infos, Sprite{
				Height:    inf.Height,
				YOffset:   inf.YOffset,
				MultiTile: inf.MultiTile,
			})

			index++
		}
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

func (s *Sprites) Get(idx int) (*ebiten.Image, *Sprite) {
	return s.sprites[idx], &s.infos[idx]
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
	if s.infos[idx].MultiTile {
		return idx + int(dirs)
	}
	return idx
}
