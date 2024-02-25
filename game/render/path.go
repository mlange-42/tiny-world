package render

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/nav"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Path is a system to render paths.
type Path struct {
	screen  generic.Resource[res.EbitenImage]
	landUse generic.Resource[res.LandUse]
	view    generic.Resource[res.View]
	aStar   nav.AStar

	path  []comp.Tile
	tiles []comp.Tile
	tick  int64
}

// InitializeUI the system
func (s *Path) InitializeUI(world *ecs.World) {
	s.screen = generic.NewResource[res.EbitenImage](world)
	s.landUse = generic.NewResource[res.LandUse](world)
	s.view = generic.NewResource[res.View](world)

	s.aStar = nav.NewAStar(s.landUse.Get())
}

// UpdateUI the system
func (s *Path) UpdateUI(world *ecs.World) {
	landUse := s.landUse.Get()

	if s.tick%60 == 0 {
		for i := 0; i < landUse.Width(); i++ {
			for j := 0; j < landUse.Height(); j++ {
				lu := landUse.Get(i, j)
				if /* lu == terr.Path || */ terr.Buildings.Contains(lu) {
					s.tiles = append(s.tiles, comp.Tile{Point: image.Pt(i, j)})
				}
			}
		}

		start := s.tiles[rand.Intn(len(s.tiles))]
		target := s.tiles[rand.Intn(len(s.tiles))]

		path, err := s.aStar.FindPath(start, target)
		if err != nil {
			fmt.Println(err)
		} else {
			s.path = path
		}
		s.tiles = s.tiles[:0]
	}
	s.tick++

	view := s.view.Get()
	off := view.Offset()
	canvas := s.screen.Get()
	img := canvas.Image

	h := view.TileHeight / 2
	z := float32(view.Zoom)
	col := color.RGBA{60, 60, 255, 255}

	n := len(s.path) - 1
	for i := 0; i < n; i++ {
		p1 := s.path[i]
		p2 := s.path[i+1]

		point1 := view.TileToGlobal(p1.X, p1.Y)
		point2 := view.TileToGlobal(p2.X, p2.Y)

		x1 := float32(point1.X)*z - float32(off.X)
		y1 := float32(point1.Y-h)*z - float32(off.Y)
		x2 := float32(point2.X)*z - float32(off.X)
		y2 := float32(point2.Y-h)*z - float32(off.Y)

		vector.StrokeLine(img, x1, y1, x2, y2, 2, col, false)
	}
}

// PostUpdateUI the system
func (s *Path) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *Path) FinalizeUI(world *ecs.World) {}
