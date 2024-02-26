package sys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// Pause system.
type Pause struct {
	speedRes generic.Resource[res.GameSpeed]

	PauseKey ebiten.Key
}

// Initialize the system
func (s *Pause) Initialize(world *ecs.World) {
	s.speedRes = generic.NewResource[res.GameSpeed](world)
}

// Update the system
func (s *Pause) Update(world *ecs.World) {
	if inpututil.IsKeyJustPressed(s.PauseKey) {
		speed := s.speedRes.Get()
		speed.Pause = !speed.Pause
	}
}

// Finalize the system
func (s *Pause) Finalize(world *ecs.World) {}
