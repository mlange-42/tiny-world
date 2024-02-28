package sys

import (
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/res"
)

// Tick system.
type Tick struct {
	speed generic.Resource[res.GameSpeed]
	time  generic.Resource[res.GameTick]
}

// Initialize the system
func (s *Tick) Initialize(world *ecs.World) {
	s.speed = generic.NewResource[res.GameSpeed](world)
	s.time = generic.NewResource[res.GameTick](world)
}

// Update the system
func (s *Tick) Update(world *ecs.World) {
	if s.speed.Get().Pause {
		return
	}

	s.time.Get().Tick++
}

// Finalize the system
func (s *Tick) Finalize(world *ecs.World) {}
