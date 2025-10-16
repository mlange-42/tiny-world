package menu

import (
	"github.com/mlange-42/ark/ecs"
	"github.com/mlange-42/tiny-world/game/res"
)

// DrawUI is a system to render the user interface.
type DrawUI struct {
	screen  ecs.Resource[res.Screen]
	ui      ecs.Resource[UI]
	sprites ecs.Resource[res.Sprites]
}

// InitializeUI the system
func (s *DrawUI) InitializeUI(world *ecs.World) {
	s.ui = s.ui.New(world)
	s.screen = s.screen.New(world)
	s.sprites = s.sprites.New(world)
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
