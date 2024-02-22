package game

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/util"
)

// EbitenImage resource for drawing.
type EbitenImage struct {
	Image  *ebiten.Image
	Width  int
	Height int
}

// Terrain resource
type Terrain struct {
	Grid[ecs.Entity]
}

type Sprites struct {
	atlas   []*ebiten.Image
	sprites []*ebiten.Image
	infos   []util.SpriteInfo
}

func NewSprites(dir string) Sprites {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	atlas := []*ebiten.Image{}
	sprites := []*ebiten.Image{}
	infos := []util.SpriteInfo{}

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

		infos := util.SpriteSheet{}
		content, err := os.ReadFile(path.Join(dir, e.Name()))
		if err != nil {
			log.Fatal("error loading JSON file: ", err)
		}
		if err := json.Unmarshal(content, &infos); err != nil {
			log.Fatal("error decoding JSON: ", err)
		}

		img, _, err := ebitenutil.NewImageFromFile(pngPath)
		if err != nil {
			log.Fatal("error reading image: ", err)
		}
		atlas = append(atlas, img)

		w, h := infos.SpriteWidth, infos.SpriteHeight
		cols, rows := img.Bounds().Dx()/w, img.Bounds().Dy()/h
		numSprites := rows * cols

		for i := 0; i < numSprites; i++ {
			row := i / cols
			col := i % cols
			sprites = append(sprites, img.SubImage(image.Rect(col*w, row*h, col*w+w, row*h+h)).(*ebiten.Image))
		}
	}

	return Sprites{
		atlas:   atlas,
		sprites: sprites,
		infos:   infos,
	}
}

func (s *Sprites) Get(idx int) *ebiten.Image {
	return s.sprites[idx]
}

type View struct {
	TileWidth, TileHeight int
	X, Y                  int
	Zoom                  float64
}

func (v *View) Offset() (int, int) {
	return int(float64(v.X) * v.Zoom), int(float64(v.Y) * v.Zoom)
}

func (v *View) Bounds(w, h int) image.Rectangle {
	vw, vh := int(float64(w)/v.Zoom), int(float64(h)/v.Zoom)

	return image.Rect(
		v.X-v.TileWidth, v.Y-3*v.TileHeight,
		v.X+vw, v.Y+vh-2*v.TileHeight,
	)
}

func (v View) TileToScreen(x, y int) image.Point {
	return image.Pt((x-y)*v.TileWidth/2,
		(x+y)*v.TileHeight/2)
}

func (v View) MouseToLocal(x, y int) (int, int) {
	return v.X + int(float64(x)/v.Zoom),
		v.Y + int(float64(y)/v.Zoom)
}
