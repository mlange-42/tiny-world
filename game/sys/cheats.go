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
	rules generic.Resource[res.Rules]
	stock generic.Resource[res.Stock]
	ui    generic.Resource[res.UI]
}

// Initialize the system
func (s *Cheats) Initialize(world *ecs.World) {
	s.rules = generic.NewResource[res.Rules](world)
	s.stock = generic.NewResource[res.Stock](world)
	s.ui = generic.NewResource[res.UI](world)
}

// Update the system
func (s *Cheats) Update(world *ecs.World) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) &&
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
		ebiten.IsKeyPressed(ebiten.KeyAlt) &&
		inpututil.IsKeyJustPressed(ebiten.KeyR) {

		stock := s.stock.Get()
		copy(stock.Res, stock.Cap)
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyShift) &&
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
		ebiten.IsKeyPressed(ebiten.KeyAlt) &&
		inpututil.IsKeyJustPressed(ebiten.KeyN) {

		ui := s.ui.Get()
		ui.ReplaceAllButtons(s.rules.Get())
	}
}

// Finalize the system
func (s *Cheats) Finalize(world *ecs.World) {}
