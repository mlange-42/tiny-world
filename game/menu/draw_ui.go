package menu

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// DrawUI is a system to render the user interface.
type DrawUI struct {
	screen  generic.Resource[res.Screen]
	ui      generic.Resource[UI]
	sprites generic.Resource[res.Sprites]
}

// InitializeUI the system
func (s *DrawUI) InitializeUI(world *ecs.World) {
	s.ui = generic.NewResource[UI](world)
	s.screen = generic.NewResource[res.Screen](world)
	s.sprites = generic.NewResource[res.Sprites](world)
}

// UpdateUI the system
func (s *DrawUI) UpdateUI(world *ecs.World) {
	screen := s.screen.Get()
	ui := s.ui.Get()

	screen.Image.Fill(s.sprites.Get().Background)

	ui.UI().Draw(screen.Image)
}

// PostUpdateUI the system
func (s *DrawUI) PostUpdateUI(world *ecs.World) {}

// FinalizeUI the system
func (s *DrawUI) FinalizeUI(world *ecs.World) {}
