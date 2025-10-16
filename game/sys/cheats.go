package sys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/res"
)

// Cheats system.
type Cheats struct {
	rules  ecs.Resource[res.Rules]
	stock  ecs.Resource[res.Stock]
	ui     ecs.Resource[res.UI]
	editor ecs.Resource[res.EditorMode]
}

// Initialize the system
func (s *Cheats) Initialize(world *ecs.World) {
	s.rules = ecs.NewResource[res.Rules](world)
	s.stock = ecs.NewResource[res.Stock](world)
	s.ui = ecs.NewResource[res.UI](world)
	s.editor = ecs.NewResource[res.EditorMode](world)
}

// Update the system
func (s *Cheats) Update(world *ecs.World) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) &&
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
		ebiten.IsKeyPressed(ebiten.KeyAlt) &&
		inpututil.IsKeyJustPressed(ebiten.KeyR) {

		if s.editor.Get().IsEditor {
			println("cheats are not available in editor mode")
			return
		}

		stock := s.stock.Get()
		copy(stock.Res, stock.Cap)
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyShift) &&
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
		ebiten.IsKeyPressed(ebiten.KeyAlt) &&
		inpututil.IsKeyJustPressed(ebiten.KeyN) {

		if s.editor.Get().IsEditor {
			println("cheats are not available in editor mode")
			return
		}

		ui := s.ui.Get()
		ui.ReplaceAllButtons(s.rules.Get())
	}
}

// Finalize the system
func (s *Cheats) Finalize(world *ecs.World) {}
