package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mlange-42/arche/ecs"
)

// EbitenImage resource for drawing.
type EbitenImage struct {
	Image  *ebiten.Image
	Width  int
	Height int
}

// Terrain resource
type Terrain = Grid[ecs.Entity]

type TerrainSprites struct {
	atlas   *ebiten.Image
	sprites []*ebiten.Image
	size    int
}

func NewTerrainSprites(file string, size int) TerrainSprites {
	img, _, err := ebitenutil.NewImageFromFile(file)
	if err != nil {
		panic(err)
	}
	cols, rows := img.Bounds().Dy()/size, img.Bounds().Dx()/size
	numSprites := rows * cols

	sprites := make([]*ebiten.Image, numSprites)

	for i := 0; i < numSprites; i++ {
		row := i / cols
		col := i % cols
		sprites[i] = img.SubImage(image.Rect(col*size, row*size, col*size+size, row*size+size)).(*ebiten.Image)
	}

	return TerrainSprites{
		atlas:   img,
		sprites: sprites,
		size:    size,
	}
}

func (s *TerrainSprites) Get(idx int) *ebiten.Image {
	return s.sprites[idx]
}

func (s *TerrainSprites) Size() int {
	return s.size
}

type View struct {
	X, Y int
	Zoom int
}

func (v *View) Bounds(sc, w, h int) (iMin, jMin, iMax, jMax int) {
	vw, vh := w/v.Zoom, h/v.Zoom

	iMin, jMin = v.X/sc, v.Y/sc
	iMax, jMax = iMin+vw, jMin+vh

	return
}

func (v *View) Offset(sc int) (int, int) {
	return v.X * v.Zoom / sc, v.Y * v.Zoom / sc
}

func (v View) MouseToLocal(sc, x, y int) (int, int) {
	return v.X + x*sc/v.Zoom,
		v.Y + y*sc/v.Zoom
}
