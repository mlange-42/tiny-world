package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// UpdateUI system.
type UpdateUI struct {
	ui generic.Resource[UI]
}

// Initialize the system
func (s *UpdateUI) Initialize(world *ecs.World) {
	s.ui = generic.NewResource[UI](world)
}

// Update the system
func (s *UpdateUI) Update(world *ecs.World) {
	ui := s.ui.Get()

	if ebiten.IsKeyPressed(ebiten.KeyShift) &&
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
		ebiten.IsKeyPressed(ebiten.KeyAlt) &&
		inpututil.IsKeyJustPressed(ebiten.KeyU) {
		ui.UnlockAll()
	}

	ui.UI().Update()
}

// Finalize the system
func (s *UpdateUI) Finalize(world *ecs.World) {}
