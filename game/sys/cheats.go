package sys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// Cheats system.
type Cheats struct {
	stock generic.Resource[res.Stock]
}

// Initialize the system
func (s *Cheats) Initialize(world *ecs.World) {
	s.stock = generic.NewResource[res.Stock](world)
}

// Update the system
func (s *Cheats) Update(world *ecs.World) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) &&
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
		inpututil.IsKeyJustPressed(ebiten.KeyR) {

		stock := s.stock.Get()
		stock.Res = stock.Cap
		return
	}
}

// Finalize the system
func (s *Cheats) Finalize(world *ecs.World) {}
