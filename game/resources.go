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
type Terrain struct {
	Grid[ecs.Entity]
}

type TerrainSprites struct {
	atlas   *ebiten.Image
	sprites []*ebiten.Image
	width   int
	height  int
}

func NewTerrainSprites(file string, width, height int) TerrainSprites {
	img, _, err := ebitenutil.NewImageFromFile(file)
	if err != nil {
		panic(err)
	}
	cols, rows := img.Bounds().Dx()/width, img.Bounds().Dy()/height
	numSprites := rows * cols

	sprites := make([]*ebiten.Image, numSprites)

	for i := 0; i < numSprites; i++ {
		row := i / cols
		col := i % cols
		sprites[i] = img.SubImage(image.Rect(col*width, row*height, col*width+width, row*height+height)).(*ebiten.Image)
	}

	return TerrainSprites{
		atlas:   img,
		sprites: sprites,
		width:   width,
		height:  height,
	}
}

func (s *TerrainSprites) Get(idx int) *ebiten.Image {
	return s.sprites[idx]
}

func (s *TerrainSprites) Size() (int, int) {
	return s.width, s.height
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
