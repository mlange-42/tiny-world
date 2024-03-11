package res

import (
	"fmt"
	"image"
	"image/color"
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
const tileSetFile = "tileset.json"

// Sprites holds all tileset data.
type Sprites struct {
	// Tile width for rendering.
	TileWidth int
	// Tile height for rendering.
	TileHeight int
	// Game background color
	Background color.RGBA
	// UI text color
	TextColor color.RGBA

	atlas       []*ebiten.Image
	sprites     []*ebiten.Image
	infos       []util.Sprite
	indices     map[string]int
	terrIndices []int
	idxUnknown  int
}

// NewSprites creates a new Sprites resource from the given tileset folder.
func NewSprites(fSys fs.FS, dir, tileSet string) Sprites {
	base := path.Join(dir, tileSet)

	tilesetJs := util.TileSet{}
	if err := util.FromJsonFs(fSys, path.Join(base, tileSetFile), &tilesetJs); err != nil {
		log.Fatal("error decoding JSON: ", err)
	}

	sheets, err := fs.ReadDir(fSys, base)
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
		if sheetFile.IsDir() || sheetFile.Name() == tileSetFile {
			continue
		}
		ext := filepath.Ext(sheetFile.Name())
		if ext != ".json" && ext != ".JSON" {
			continue
		}
		baseName := strings.Replace(sheetFile.Name(), ext, "", 1)
		pngPath := path.Join(base, fmt.Sprintf("%s.png", baseName))

		sheet := util.SpriteSheet{}
		if err := util.FromJsonFs(fSys, path.Join(base, sheetFile.Name()), &sheet); err != nil {
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

			if inf.AnimSpeed == 0 {
				inf.AnimSpeed = 1
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

	terrIndices := make([]int, len(terr.Properties))
	for i := range terr.Properties {
		if idx, ok := indices[terr.Properties[i].Name]; ok {
			terrIndices[i] = idx
		} else {
			terrIndices[i] = indices[nameUnknown]
		}
	}

	return Sprites{
		TileWidth:   tilesetJs.TileWidth,
		TileHeight:  tilesetJs.TileHeight,
		Background:  tilesetJs.BackgroundColor,
		TextColor:   tilesetJs.TextColor,
		atlas:       atlas,
		sprites:     sprites,
		infos:       infos,
		indices:     indices,
		idxUnknown:  indices[nameUnknown],
		terrIndices: terrIndices,
	}
}

// GetInfo returns the sprite info for an index.
func (s *Sprites) GetInfo(idx int) *util.Sprite {
	return &s.infos[idx]
}

// Get returns the sprite image for an index.
func (s *Sprites) Get(idx int) *ebiten.Image {
	inf := &s.infos[idx]
	return s.sprites[inf.Index[0]]
}

// GetRand returns the sprite image for an index, with animation and random variations.
func (s *Sprites) GetRand(idx int, frame int, rand int) *ebiten.Image {
	inf := &s.infos[idx]

	if inf.IsAnimated() {
		vars := len(inf.Index) / inf.AnimFrames
		sIdx := (rand%vars)*inf.AnimFrames + (frame/inf.AnimSpeed)%inf.AnimFrames
		return s.sprites[inf.Index[sIdx]]
	} else {
		return s.sprites[inf.Index[rand%len(inf.Index)]]
	}
}

// GetSprite returns the sprite image for an index, including sub-index.
func (s *Sprites) GetSprite(idx int) *ebiten.Image {
	return s.sprites[idx]
}

// GetIndex returns the sprite index for a sprite or terrain ID.
func (s *Sprites) GetIndex(name string) int {
	if idx, ok := s.indices[name]; ok {
		return idx
	}
	return s.idxUnknown
}

// GetTerrainIndex returns the sprite index for a terrain ID.
func (s *Sprites) GetTerrainIndex(t terr.Terrain) int {
	return s.terrIndices[t]
}

// GetMultiTileTerrainIndex returns the sprite index for a terrain ID, using multitile.
func (s *Sprites) GetMultiTileTerrainIndex(t terr.Terrain, dirs terr.Directions, frame int, rand int) int {
	idx := s.terrIndices[t]
	return s.GetMultiTileIndex(idx, dirs, frame, rand)
}

// GetMultiTileTerrainIndex returns the sprite index for an index, using multitile.
func (s *Sprites) GetMultiTileIndex(idx int, dirs terr.Directions, frame int, rand int) int {
	inf := &s.infos[idx]
	if inf.IsMultitile() {
		sprites := inf.Multitile[dirs]
		if inf.IsAnimated() {
			vars := len(sprites) / inf.AnimFrames
			sIdx := (rand%vars)*inf.AnimFrames + (frame/inf.AnimSpeed)%inf.AnimFrames
			return sprites[sIdx]
		} else {
			return sprites[rand%len(sprites)]
		}
	} else if inf.IsAnimated() {
		vars := len(inf.Index) / inf.AnimFrames
		sIdx := (rand%vars)*inf.AnimFrames + (frame/inf.AnimSpeed)%inf.AnimFrames
		return inf.Index[sIdx]
	}
	return idx
}
