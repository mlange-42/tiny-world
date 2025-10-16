package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
)

// HaulerPaths is a system to render paths.
type HaulerPaths struct {
	screen ecs.Resource[res.Screen]
	view   ecs.Resource[res.View]
	update ecs.Resource[res.UpdateInterval]

	filter *ecs.Filter1[comp.Hauler]
}

// InitializeUI the system
func (s *HaulerPaths) InitializeUI(world *ecs.World) {
	s.screen = s.screen.New(world)
	s.view = s.view.New(world)
	s.update = s.update.New(world)

	s.filter = s.filter.New(world)
}

// UpdateUI the system
func (s *HaulerPaths) UpdateUI(world *ecs.World) {
	view := s.view.Get()
	update := s.update.Get()
	off := view.Offset()
	canvas := s.screen.Get()
	img := canvas.Image

	h := view.TileHeight / 2
	z := float32(view.Zoom)
	col := color.RGBA{60, 60, 255, 255}

	query := s.filter.Query()
	for query.Next() {
		haul := query.Get()
		path := haul.Path
		n := len(path) - 1
		for i := 0; i < n; i++ {
			p1 := path[i]
			p2 := path[i+1]

			point1 := view.TileToGlobal(p1.X, p1.Y)
			point2 := view.TileToGlobal(p2.X, p2.Y)

			x1 := float32(point1.X)*z - float32(off.X)
			y1 := float32(point1.Y-h)*z - float32(off.Y)
			x2 := float32(point2.X)*z - float32(off.X)
			y2 := float32(point2.Y-h)*z - float32(off.Y)

			//vector.StrokeLine(img, x1, y1, x2, y2, 2, col, false)

			if i != n-1 {
				continue
			}

			dt := float32(haul.PathFraction) / float32(update.Interval)
			xx := x1*dt + x2*(1-dt)
			yy := y1*dt + y2*(1-dt)

			vector.DrawFilledCircle(img, xx, yy, 4, col, false)
		}
	}
}

// PostUpdateUI the system
func (s *HaulerPaths) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *HaulerPaths) FinalizeUI(world *ecs.World) {}
